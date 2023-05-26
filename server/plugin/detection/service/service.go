package service

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	local "github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/perception"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DetectionService struct{}

func (e *DetectionService) PlugService(req model.Request) (res model.Response, err error) {
	// 写你的业务逻辑
	return res, nil
}

//@author: [neophack](https://github.com/neophack)
//@function: Upload
//@description: 创建文件上传记录
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *DetectionService) Upload(file model.DetectionFileUploadAndDownload) error {
	return global.GVA_DB.Create(&file).Error
}

//@author: [neophack](https://github.com/neophack)
//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: model.ExaFileUploadAndDownload, error

func (e *DetectionService) FindFile(id uint) (model.DetectionFileUploadAndDownload, error) {
	var file model.DetectionFileUploadAndDownload
	err := global.GVA_DB.Where("id = ?", id).First(&file).Error
	return file, err
}

//@author: [neophack](https://github.com/neophack)
//@function: DeleteFile
//@description: 删除文件记录
//@param: file model.ExaFileUploadAndDownload
//@return: err error

func (e *DetectionService) DeleteFile(file model.DetectionFileUploadAndDownload) (err error) {
	var fileFromDb model.DetectionFileUploadAndDownload
	fileFromDb, err = e.FindFile(file.ID)
	if err != nil {
		return
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}

	err = global.GVA_DB.Where("id = ?", file.ID).Unscoped().Delete(&file).Error

	return err
}

// EditFileName 编辑文件名或者备注
func (e *DetectionService) EditFileName(file model.DetectionFileUploadAndDownload) (err error) {
	var fileFromDb model.DetectionFileUploadAndDownload
	return global.GVA_DB.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

//@author: [neophack](https://github.com/neophack)
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: list interface{}, total int64, err error

func (e *DetectionService) GetFileRecordInfoList(info request.PageInfo) (fileLists []model.DetectionFileUploadAndDownload, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	//user := info.User
	//app := info.App
	batchid := info.Batchid
	db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
	//var fileLists []model.DetectionFileUploadAndDownload
	if len(batchid) > 0 {
		db = db.Where("batchid = '" + batchid + "'")
	}
	//if len(user) > 0 {
	//	db = db.Where("user = '" + user + "'")
	//}
	//if len(app) > 0 {
	//	db = db.Where("app = '" + app + "'")
	//}
	if len(keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		global.GVA_DB.AutoMigrate(model.DetectionFileUploadAndDownload{})
		return
	}
	err = db.Limit(limit).Offset(offset).Order("name asc").Find(&fileLists).Error
	return fileLists, total, err
}

//@author: [neophack](https://github.com/neophack)
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: file model.ExaFileUploadAndDownload, err error

func (e *DetectionService) UploadFile(header *multipart.FileHeader, noSave string, user string, app string, batchid string, size string) (file model.DetectionFileUploadAndDownload, err error) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	mkdirErr := os.MkdirAll(global.GVA_CONFIG.Local.TmpPath, os.ModePerm)
	if mkdirErr != nil {
		global.GVA_LOG.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
	}
	if uploadErr != nil {
		return
	}
	if noSave == "0" {

		re := regexp.MustCompile(`filename="([^"]*)"`)
		matches := re.FindAllStringSubmatch(header.Header.Get("Content-Disposition"), -1)
		filename := header.Filename
		for _, match := range matches {
			filename = match[1]
		}
		s := strings.Split(header.Filename, ".")
		f := model.DetectionFileUploadAndDownload{
			Url:     filePath,
			Name:    filename,
			Tag:     s[len(s)-1],
			Key:     key,
			Batchid: batchid,
			Size:    size,
		}
		return f, e.Upload(f)
	}
	return
}

//@author: [neophack](https://github.com/neophack)
//@function: GetBatchInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: list interface{}, total int64, err error

func (e *DetectionService) GetBatchInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	user := info.User
	app := info.App
	db := global.GVA_DB.Model(&model.DetectionFileBatch{})
	var fileLists []model.DetectionFileBatch
	if len(user) > 0 {
		db = db.Where("own = '" + user + "'")
	}
	if len(app) > 0 {
		db = db.Where("app = '" + app + "'")
	}
	if len(keyword) > 0 {
		db = db.Where("batchid LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		global.GVA_DB.AutoMigrate(model.DetectionFileBatch{})
		global.GVA_DB.AutoMigrate(model.DetectionFileUploadAndDownload{})
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&fileLists).Error
	return fileLists, total, err
}

// NewBatch 函数创建一个新的文件批次，包括用户、应用、批次ID、文件数量和文件大小，返回一个 DetectionFileBatch 结构体，并将其添加到数据库中。
func (e *DetectionService) NewBatch(user string, app string, batchid string, filesCount int, filesSize string) (file model.DetectionFileBatch, err error) {
	f := model.DetectionFileBatch{
		Batchid:    batchid,     // 批次ID
		Own:        user,        // 用户名
		App:        app,         // 应用名
		FilesCount: filesCount,  // 文件数量
		FilesSize:  filesSize,   // 文件大小
		Status:     "uploading", // 批次状态
	}

	return f, global.GVA_DB.Create(&f).Error
}

// DeleteBatch 函数删除给定批次ID、应用和状态的文件批次及其关联的文件。它首先从数据库中删除文件批次，然后删除关联的文件，最后从数据库中删除这些文件。
func (e *DetectionService) DeleteBatch(user string, app string, batchid string, status string) (err error) {
	db := global.GVA_DB.Model(&model.DetectionFileBatch{})
	err = db.Where("batchid = ?", batchid).Unscoped().Delete(&model.DetectionFileBatch{}, "app = ?", app).Error // 从数据库中删除文件批次
	if err != nil {
		return err
	}
	db2 := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
	var fileLists []model.DetectionFileUploadAndDownload
	err = db2.Where("batchid = '" + batchid + "'").Find(&fileLists).Error // 查找关联的文件
	for i := range fileLists {
		var fileFromDb model.DetectionFileUploadAndDownload
		fileFromDb, err = e.FindFile(fileLists[i].ID)
		if err != nil {
			return err
		}
		oss := upload.NewOss()
		if err = oss.DeleteFile(fileFromDb.Key); err != nil { // 删除文件
			return errors.New("文件删除失败")
		}
		if fileLists[i].UrlDetection != "" {
			err = os.Remove(fileLists[i].UrlDetection) // 删除文件
			if err != nil {
				return err
			}
		}

		err = global.GVA_DB.Where("id = ?", fileLists[i].ID).Unscoped().Delete(&fileLists[i]).Error // 从数据库中删除文件
	}

	return err
}

// formatSize 函数将给定的文件大小转换为易读的格式（例如，将字节转换为KB、MB、GB等）。
func formatSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	var i int
	floatSize := float64(size)
	for i = 0; floatSize >= 1024 && i < len(units)-1; i++ {
		floatSize /= 1024
	}
	return fmt.Sprintf("%.2f %s", floatSize, units[i])
}

// ChangeStatus 函数更新文件批次的状态，并计算文件总数和大小。
func (e *DetectionService) ChangeStatus(user string, app string, batchid string, status string) (err error) {
	db := global.GVA_DB.Model(&model.DetectionFileBatch{})
	db = db.Where("batchid = ?", batchid).Where("app = ?", app)
	err = db.Update("status", status).Error // 更新批次状态
	if err != nil {
		return err
	}
	db2 := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})

	var total int64
	db2 = db2.Where("batchid = '" + batchid + "'")
	err = db2.Count(&total).Error // 计算文件总数
	if err != nil {
		return err
	}
	var fileLists []model.DetectionFileUploadAndDownload
	var fileSize int64
	err = db2.Find(&fileLists).Error // 查找文件列表
	for i := range fileLists {
		number_int, err := strconv.Atoi(fileLists[i].Size)
		if err != nil {
			return err
		}
		fileSize += int64(number_int) // 计算文件大小
	}
	if err != nil {
		return err
	}
	err = db.Update("FilesCount", total).Error // 更新文件总数
	if err != nil {
		return err
	}
	formattedSize := formatSize(fileSize)
	fmt.Println(formattedSize)
	err = db.Update("FilesSize", formattedSize).Error // 更新文件大小
	return err
}

// Dojob 函数处理文件批次的任务。它从数据库中查找状态为“ready”的文件批次，然后启动处理程序对这些文件进行处理。
func (e *DetectionService) Dojob() {
	maxJobs := int64(3) // 设置最大任务数
	for {
		if global.GVA_DB == nil {
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
	for {
		time.Sleep(time.Second * 5)
		db := global.GVA_DB.Model(&model.DetectionFileBatch{})
		var batchLists []model.DetectionFileBatch
		db = db.Where("status = 'ready'")
		err := db.Order("created_at desc").Find(&batchLists).Error // 查找状态为“ready”的文件批次
		if err != nil {
			global.GVA_DB.AutoMigrate(model.DetectionFileBatch{})
			global.GVA_DB.AutoMigrate(model.DetectionFileUploadAndDownload{})
			continue
		}
		var workingCount int64

		for ii := range batchLists {
			app := batchLists[ii].App
			batchid := batchLists[ii].Batchid
			id := batchLists[ii].ID

			for i := range local.GlobalConfig_.ModelConfig {
				c := local.GlobalConfig_.ModelConfig[i]
				if c.App == app {
					db = global.GVA_DB.Model(&model.DetectionFileBatch{})
					err = db.Where("status = 'working'").Count(&workingCount).Error // 查找状态为“working”的文件批次
					if err != nil {
						continue
					}
					if workingCount < maxJobs {
						var isworking = false
						for j := range local.WorkingBatchs_ {
							if batchid == local.WorkingBatchs_[j] {
								isworking = true
							}
						}
						if !isworking {
							local.WorkingBatchs_ = append(local.WorkingBatchs_, batchid)
							go perception.RunBatch(c.ProgramPath, batchid, id) // 启动处理程序对文件进行处理
						}
						time.Sleep(time.Second * 1)
						i = 0
					}
				}
			}
		}
	}
}
