package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/api"
	"github.com/gin-gonic/gin"
)

type DetectionRouter struct {
}

// InitDetectionRouter 函数用于初始化检测模块的路由
// Router 参数是 gin 框架的路由组
func (s *DetectionRouter) InitDetectionRouter(Router *gin.RouterGroup) {
	// 定义一个不需要认证的路由组
	plugRouterFree := Router.Use()
	// 获取检测模块的 API
	plugApi := api.ApiGroupApp.DetectionApi

	// 定义不需要认证的路由
	{
		// 上传文件
		plugRouterFree.POST("upload", plugApi.UploadFile)
		// 下载文件列表的压缩包
		plugRouterFree.GET("downloadFilesZip", plugApi.DownloadFilesZip)
	}

	// 定义需要认证的路由组
	plugRouter := Router.Use(middleware.OperationRecord())
	{
		// 获取上传文件列表
		plugRouter.POST("getFileList", plugApi.GetFileList)
		// 删除指定文件
		plugRouter.POST("deleteFile", plugApi.DeleteFile)
		// 编辑文件名或者备注
		plugRouter.POST("editFileName", plugApi.EditFileName)
		// 获取文件上传的批次列表
		plugRouter.POST("getBatchList", plugApi.GetBatchList)
		// 新建一个文件上传批次
		plugRouter.POST("newBatch", plugApi.NewBatch)
		// 更改文件上传批次的状态
		plugRouter.POST("changeStatus", plugApi.ChangeStatus)
		// 删除文件上传批次
		plugRouter.POST("deleteBatch", plugApi.DeleteBatch)
	}
}
