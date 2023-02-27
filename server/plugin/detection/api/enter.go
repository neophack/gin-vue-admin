package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/service"

type ApiGroup struct {
	DetectionApi
}

var ApiGroupApp = new(ApiGroup)

var DetectionService = service.ServiceGroupApp.DetectionService
