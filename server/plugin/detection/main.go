package detection

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/config"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/router"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/service"
	"github.com/gin-gonic/gin"
)

type DetectionPlugin struct {
}
var DetectionService = service.ServiceGroupApp.DetectionService
func CreateDetectionPlug(config []config.ModelConfig, ) *DetectionPlugin {
	global.GlobalConfig_.ModelConfig = config
	DetectionService.Dojob()
	return &DetectionPlugin{}
}

func (*DetectionPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitDetectionRouter(group)
}

func (*DetectionPlugin) RouterPath() string {
	return "detection"
}
