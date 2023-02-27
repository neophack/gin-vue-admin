package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	local "github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/perception"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"mime/multipart"
	"strings"
)

type DetectionService struct{}

func (e *DetectionService) PlugService(req model.Request) (res model.Response, err error) {
	// 写你的业务逻辑
	return res, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Upload
//@description: 创建文件上传记录
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *DetectionService) Upload(file model.DetectionFileUploadAndDownload) error {
	return global.GVA_DB.Create(&file).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: FindFile
//@description: 查询文件记录
//@param: id uint
//@return: model.ExaFileUploadAndDownload, error

func (e *DetectionService) FindFile(id uint) (model.DetectionFileUploadAndDownload, error) {
	var file model.DetectionFileUploadAndDownload
	err := global.GVA_DB.Where("id = ?", id).First(&file).Error
	return file, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFile
//@description: 删除文件记录
//@param: file model.ExaFileUploadAndDownload
//@return: err error

func (e *DetectionService) DeleteFile(file model.DetectionFileUploadAndDownload) (err error) {
	var fileFromDb model.DetectionFileUploadAndDownload
	//fileFromDb, err = e.FindFile(file.ID)
	//if err != nil {
	//	return
	//}
	//oss := upload.NewOss()
	//if err = oss.DeleteFile(fileFromDb.Key); err != nil {
	//	return errors.New("文件删除失败")
	//}

	db := global.GVA_DB.Model(&system.SysUser{})
	var user system.SysUser
	err = db.Where("id = 1").First(&user).Error
	if err != nil {
		return
	}
	if user.UUID.String() == file.Own {
		err = global.GVA_DB.Where("id = ?", file.ID).Unscoped().Delete(&file).Error
	} else {
		err = global.GVA_DB.Where("id = ?", file.ID).First(&fileFromDb).Update("own", user.UUID.String()).Error
	}

	return err
}

// EditFileName 编辑文件名或者备注
func (e *DetectionService) EditFileName(file model.DetectionFileUploadAndDownload) (err error) {
	var fileFromDb model.DetectionFileUploadAndDownload
	return global.GVA_DB.Where("id = ?", file.ID).First(&fileFromDb).Update("name", file.Name).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: list interface{}, total int64, err error

func (e *DetectionService) GetFileRecordInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Keyword
	user := info.User
	app := info.App
	db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
	var fileLists []model.DetectionFileUploadAndDownload
	if len(user) > 0 {
		db = db.Where("own = '" + user + "'")
	}
	if len(app) > 0 {
		db = db.Where("app = '" + app + "'")
	}
	if len(keyword) > 0 {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&fileLists).Error
	return fileLists, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: file model.ExaFileUploadAndDownload, err error

func (e *DetectionService) UploadFile(header *multipart.FileHeader, noSave string, user string, app string) (file model.DetectionFileUploadAndDownload, err error) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(err)
	}
	if noSave == "0" {
		s := strings.Split(header.Filename, ".")
		f := model.DetectionFileUploadAndDownload{
			Url:  filePath,
			Name: header.Filename,
			Tag:  s[len(s)-1],
			Key:  key,
			Own:  user,
			App:  app,
		}
		return f, e.Upload(f)
	}
	return
}

func (e *DetectionService) Dojob() {
	for i := range local.GlobalConfig_.ModelConfig {
		c := local.GlobalConfig_.ModelConfig[i]
		if c.Algorithm == "yolov8seg" {
			go perception.Yolov8seg(c.ModelPath, c.App)
		} else if c.Algorithm == "yolov5" {
			go perception.Yolov5(c.ModelPath, c.App)
		}
	}
	//perception.Yolov5("yolov5s.onnx")
	//perception.Yolov8seg("yolov5s.onnx")
}
