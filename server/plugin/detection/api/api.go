package api

import (
	"archive/zip"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
)

type DetectionApi struct{}

// UploadFile
// @Tags detection
// @Summary 上传文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param   file   formData                                             true "上传文件"
// @Param   noSave query     string                                     false "是否保存文件到本地"
// @Success 200    {object}  response.Response{msg=string}              "上传文件示例,返回接收文件失败&上传成功"
// @Router /detection/uploadFile [post]
func (b *DetectionApi) UploadFile(c *gin.Context) {
	var file model.DetectionFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	userName := c.Request.Header.Get("User")
	appName := c.Request.Header.Get("App")
	batchid := c.Request.Header.Get("Batchid")
	size := c.Request.Header.Get("Content-Length")
	//fmt.Println("size:", size)
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file, err = DetectionService.UploadFile(header, noSave, userName, appName, batchid, size) // 文件上传后拿到文件路径
	if err != nil {
		global.GVA_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	//if progress > 99.99 {
	//	db := global.GVA_DB.Model(&model.DetectionFileBatch{})
	//	err = db.Where("batchid = ?", batchid).Update("status", "finish").Error
	//	if err != nil {
	//		return
	//	}
	//}

	response.OkWithDetailed(model.DetectionFileResponse{File: file}, "上传成功", c)
}

// EditFileName
// @Tags detection
// @Summary 编辑文件名或者备注
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param   file model.DetectionFileUploadAndDownload true "编辑文件名或者备注"
// @Success 200 {object} response.Response{msg=string} "编辑文件名或者备注,返回编辑失败&编辑成功"
// @Router /detection/editFileName [post]
func (b *DetectionApi) EditFileName(c *gin.Context) {
	var file model.DetectionFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = DetectionService.EditFileName(file)
	if err != nil {
		global.GVA_LOG.Error("编辑失败!", zap.Error(err))
		response.FailWithMessage("编辑失败", c)
		return
	}
	response.OkWithMessage("编辑成功", c)
}

// DeleteFile
// @Tags detection
// @Summary 删除文件
// @Security ApiKeyAuth
// @Produce application/json
// @Param data body model.ExaFileUploadAndDownload true "传入文件里面id即可"
// @Success 200 {object} response.Response{msg=string} "删除文件"
// @Router /detection/deleteFile [post]
func (b *DetectionApi) DeleteFile(c *gin.Context) {
	var file model.DetectionFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := DetectionService.DeleteFile(file); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetFileList
// @Tags detection
// @Summary 分页文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页文件列表,返回包括列表,总数,页码,每页数量
// @Router /detection/getFileList [post]
func (b *DetectionApi) GetFileList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := DetectionService.GetFileRecordInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetBatchList
// @Tags detection
// @Summary 分页获取文件批次列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取文件批次列表,返回包括列表,总数,页码,每页数量"
// @Router /detection/getBatchList [post]
func (b *DetectionApi) GetBatchList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := DetectionService.GetBatchInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// NewBatch
// @Tags detection
// @Summary 新建文件批次
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DetectionFileBatch true "DetectionFileBatch"
// @Success 200 {object} response.Response "新建文件批次"
// @Router /detection/newBatch [post]
func (b *DetectionApi) NewBatch(c *gin.Context) {
	var pageInfo model.DetectionFileBatch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userName := pageInfo.Own
	appName := pageInfo.App
	batchName := pageInfo.Batchid
	filesCount := pageInfo.FilesCount
	filesSize := pageInfo.FilesSize

	file, err := DetectionService.NewBatch(userName, appName, batchName, filesCount, filesSize) // 文件上传后拿到文件路径
	if err != nil {
		global.GVA_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(model.DetectionBatchResponse{File: file}, "new成功", c)
}

// ChangeStatus
// @Tags detection
// @Summary 修改文件批次状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DetectionFileBatch true "DetectionFileBatch"
// @Success 200 {object} response.Response "修改文件批次状态"
// @Router /detection/changeStatus [post]
func (b *DetectionApi) ChangeStatus(c *gin.Context) {
	var pageInfo model.DetectionFileBatch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userName := pageInfo.Own
	appName := pageInfo.App
	batchName := pageInfo.Batchid
	status := pageInfo.Status

	err = DetectionService.ChangeStatus(userName, appName, batchName, status) // 文件上传后拿到文件路径
	if err != nil {
		global.GVA_LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"data": "ok"}, "update成功", c)
}

// DeleteBatch
// @Tags detection
// @Summary 删除文件批次
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DetectionFileBatch true "DetectionFileBatch"
// @Success 200 {object} response.Response "删除文件批次"
// @Router /detection/deleteBatch [post]
func (b *DetectionApi) DeleteBatch(c *gin.Context) {
	var pageInfo model.DetectionFileBatch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userName := pageInfo.Own
	appName := pageInfo.App
	batchName := pageInfo.Batchid
	status := pageInfo.Status

	err = DetectionService.DeleteBatch(userName, appName, batchName, status)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// DownloadFilesZip
// @Tags detection
// @Summary 批量下载文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param batchid query string true "批次id"
// @Param result query string true "结果类型"
// @Success 200 {object} response.Response "批量下载文件"
// @Router /detection/downloadFilesZip [get]
func (b *DetectionApi) DownloadFilesZip(c *gin.Context) {
	// 获取要下载的文件列表
	batchid := c.Query("batchid")
	result := c.Query("result")
	if batchid == "" {
		response.FailWithMessage("获取失败", c)
		return
	}
	fmt.Println("batchid:", batchid)
	var pageInfo request.PageInfo
	pageInfo.Batchid = batchid
	//err := c.ShouldBindJSON(&pageInfo)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	flist, _, err := DetectionService.GetFileRecordInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	// 设置响应头部信息
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+batchid+".zip")
	c.Header("Content-Transfer-Encoding", "binary")

	// 创建一个带缓冲的管道
	//buf := new(bytes.Buffer)
	writer := zip.NewWriter(c.Writer)

	// 遍历要压缩的文件列表，并将文件内容写入管道中
	for ii := range flist {
		file := flist[ii].Url
		if result != "" {
			file = flist[ii].UrlDetection
		}
		f, err := os.Open(file)
		if err != nil {
			log.Printf("Failed to open file %s: %v\n", file, err)
			continue
		}

		// 创建一个压缩文件的头部信息
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("Failed to get file info for %s: %v\n", file, err)
			continue
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Printf("Failed to create header for %s: %v\n", file, err)
			continue
		}

		header.Name = flist[ii].Name

		// 将文件内容写入管道中
		fw, err := writer.CreateHeader(header)
		if err != nil {
			log.Printf("Failed to create writer for %s: %v\n", file, err)
			continue
		}

		_, err = io.Copy(fw, f)
		if err != nil {
			log.Printf("Failed to write file content to zip: %v\n", err)
			continue
		}

		f.Close()
	}

	// 关闭压缩写入器
	err = writer.Close()
	if err != nil {
		log.Printf("Failed to close zip writer: %v\n", err)
		return
	}

	//// 将压缩数据写入响应体中
	//_, err = io.Copy(c.Writer, buf)
	//if err != nil {
	//	log.Printf("Failed to write zip data to response: %v\n", err)
	//	return
	//}
}
