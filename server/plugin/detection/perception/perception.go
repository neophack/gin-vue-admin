package perception

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	local "github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/model"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// RunBatch 函数用于启动一个文件批次的检测任务
// programPath 是检测程序的路径，batchid 是文件批次的 ID，id 是文件批次在数据库中的记录 ID
func RunBatch(programPath string, batchid string, id uint) {
	// 更新文件批次的状态为 working
	db2 := global.GVA_DB.Model(&model.DetectionFileBatch{})
	err := db2.Where("id = ?", id).Update("status", "working").Error
	if err != nil {
		return
	}

	// 创建临时目录，用于存放检测结果
	mkdirErr := os.MkdirAll(global.GVA_CONFIG.Local.TmpPath, os.ModePerm)
	if mkdirErr != nil {
		global.GVA_LOG.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
	}

	// 从数据库中获取当前文件批次的所有文件列表
	db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
	var fileLists []model.DetectionFileUploadAndDownload
	db = db.Where("batchid = '" + batchid + "'")
	err = db.Order("created_at desc").Find(&fileLists).Error
	if err != nil {
		return
	}

	// 生成一个随机的文件名，用于存放文件列表
	rand.Seed(time.Now().UnixNano())
	fileSuffix := fmt.Sprintf("%d", rand.Intn(1000))
	fileName := "detlist-" + fileSuffix + ".txt"

	// 创建临时文件，注意，不需要在文件名中指定路径，ioutil.TempFile 会自动使用系统默认的临时文件夹
	f, err := ioutil.TempFile("", fileName)
	if err != nil {
		log.Fatal(err)
	}

	// 使用 bufio.NewWriter 函数创建一个带缓冲的写入器
	w := bufio.NewWriter(f)

	// 将文件列表写入临时文件
	for ii := range fileLists {
		imgPath := fileLists[ii].Url
		_, err := w.WriteString(imgPath + "\n")
		if err != nil {
			log.Fatal(err)
		}

		// 如果当前文件已经检测过了，删除之前的检测结果文件
		if fileLists[ii].UrlDetection != "" {
			_, err = os.Stat(fileLists[ii].UrlDetection)
			if !os.IsNotExist(err) {
				err = os.Remove(fileLists[ii].UrlDetection)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}
	}

	// 调用 Flush 方法，将缓冲中的数据写入文件
	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}

	// 关闭文件
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 打印生成的文件名
	fmt.Println("The file", f.Name(), "has been created.")

	// 启动检测程序，并传递命令行参数
	cmd := exec.Command(programPath, "-i", f.Name(), "-o", global.GVA_CONFIG.Local.TmpPath, "--gpu", "--imlist")
	output, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// 获取检测程序的进程 ID
	pid := cmd.Process.Pid
	fmt.Println("PID: ", pid)

	// 启动一个 goroutine，用于读取检测程序的输出，并更新文件批次的进度和状态
	go func() {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Progress: ") {
				// 解析检测进度
				progressStr := strings.TrimPrefix(line, "Progress: ")
				progress, err := strconv.ParseFloat(progressStr, 8)

				if err != nil {
					panic(err)
				}

				// 更新文件批次的进度
				fmt.Printf("%s Progress: %.1f%%\r", batchid, progress)
				err = db2.Update("progress", fmt.Sprintf("%.1f", progress)).Error
				if err != nil {
					return
				}

				// 如果检测完成，更新文件批次的状态，并将检测结果存入数据库
				if progress > 99.99 {
					err = db2.Update("status", "finish").Error
					if err != nil {
						return
					}

					for i := range fileLists {
						newurl := strings.Replace(fileLists[i].Url, global.GVA_CONFIG.Local.StorePath, global.GVA_CONFIG.Local.TmpPath, 1)
						_, err = os.Stat(newurl)
						if !os.IsNotExist(err) {
							db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
							err = db.Where("id = ?", fileLists[i].ID).Update("url_detection", newurl).Error
							if err != nil {
								return
							}
						}
					}
				} else {
					// 如果检测未完成，更新文件批次的状态为 working
					err = db2.Update("status", "working").Error
					if err != nil {
						return
					}
				}
			}
		}
	}()

	// 等待检测程序退出，并处理检测程序的错误
	if err := cmd.Wait(); err != nil {
		fmt.Println("exit error!")
		err = db2.Update("status", "error").Error
		if err != nil {
			return
		}
		return
	}
	fmt.Println("finish!")

	// 从正在运行的文件批次列表中删除当前文件批次
	newSlice := []string{}
	for _, s := range local.WorkingBatchs_ {
		if s != batchid {
			newSlice = append(newSlice, s)
		}
	}
	local.WorkingBatchs_ = newSlice
}
