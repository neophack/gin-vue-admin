package perception

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
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

func RunBatch(programPath string, batchid string, id uint) {
	db2 := global.GVA_DB.Model(&model.DetectionFileBatch{})
	err := db2.Where("id = ?", id).Update("status", "working").Error
	if err != nil {
		return
	}

	mkdirErr := os.MkdirAll(global.GVA_CONFIG.Local.TmpPath, os.ModePerm)
	if mkdirErr != nil {
		global.GVA_LOG.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
	}
	db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
	var fileLists []model.DetectionFileUploadAndDownload
	db = db.Where("batchid = '" + batchid + "'")
	//db = db.Where("url_detection = '' or url_detection isnull")
	err = db.Order("created_at desc").Find(&fileLists).Error
	if err != nil {
		return
	}

	// 生成一个随机的文件名
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

	// 将字符串列表写入文件
	for ii := range fileLists {
		imgPath := fileLists[ii].Url
		_, err := w.WriteString(imgPath + "\n")
		if err != nil {
			log.Fatal(err)
		}
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

	cmd := exec.Command(programPath, "-i", f.Name(), "-o", global.GVA_CONFIG.Local.TmpPath, "--gpu", "--imlist")
	output, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	// get pid of the process
	pid := cmd.Process.Pid
	fmt.Println("PID: ", pid)

	go func() {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.HasPrefix(line, "Progress: ") {
				progressStr := strings.TrimPrefix(line, "Progress: ")
				progress, err := strconv.ParseFloat(progressStr, 8)

				if err != nil {
					panic(err)
				}
				fmt.Printf("Progress: %.1f%%\n", progress)
				//fmt.Println(progressStr)

				err = db2.Where("id = ?", id).Update("progress", fmt.Sprintf("%.1f", progress)).Error
				if err != nil {
					return
				}
				if progress > 99.99 {
					err = db2.Where("id = ?", id).Update("status", "finish").Error
					if err != nil {
						return
					}

					for i := range fileLists {
						newurl := strings.Replace(fileLists[i].Url, global.GVA_CONFIG.Local.StorePath, global.GVA_CONFIG.Local.TmpPath, 1)
						_, err = os.Stat(newurl)
						fmt.Println("IsNotExist", os.IsNotExist(err))
						if !os.IsNotExist(err) {
							db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
							err = db.Where("id = ?", fileLists[i].ID).Update("url_detection", newurl).Error
							fmt.Println(fileLists[i].ID, newurl)
							if err != nil {
								return
							}
						}

					}
				}

			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("exit error!")
		err = db2.Where("id = ?", id).Update("status", "error").Error
		if err != nil {
			return
		}
		return
	}
	fmt.Println("finish!")
}
