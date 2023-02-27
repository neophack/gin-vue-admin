package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/api"
	"github.com/gin-gonic/gin"
)

type DetectionRouter struct {
}

func (s *DetectionRouter) InitDetectionRouter(Router *gin.RouterGroup) {
	plugRouterFree := Router.Use()
	plugApi := api.ApiGroupApp.DetectionApi
	{
		plugRouterFree.POST("upload", plugApi.UploadFile)                                 // 上传文件
	}
	plugRouter := Router.Use(middleware.OperationRecord())
	{
		//plugRouter.POST("routerName", plugApi.ApiName)
		//plugRouterFree.POST("upload", plugApi.UploadFile)                                 // 上传文件
		plugRouter.POST("getFileList", plugApi.GetFileList)                           // 获取上传文件列表
		plugRouter.POST("deleteFile", plugApi.DeleteFile)                             // 删除指定文件
		plugRouter.POST("editFileName", plugApi.EditFileName)                         // 编辑文件名或者备注
		//plugRouter.POST("breakpointContinue", plugApi.BreakpointContinue)             // 断点续传
		//plugRouter.GET("findFile", plugApi.FindFile)                                  // 查询当前文件成功的切片
		//plugRouter.POST("breakpointContinueFinish", plugApi.BreakpointContinueFinish) // 切片传输完成
		//plugRouter.POST("removeChunk", plugApi.RemoveChunk)                           // 删除切片
	}
}
