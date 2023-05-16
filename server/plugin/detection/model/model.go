package model

import (
	"gorm.io/gorm"
	"time"
)

type Request struct {
	r_1 string // r_1
}

type Response struct {
	rp_1 string // rp_1
}
type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

type DetectionFileUploadAndDownload struct {
	GVA_MODEL
	Name         string `json:"name" gorm:"comment:文件名"` // 文件名
	Batchid      string `json:"batchid" gorm:"comment:batchid"`
	Url          string `json:"url" gorm:"comment:文件地址"`                    // 文件地址
	UrlDetection string `json:"url_detection" gorm:"comment:detection文件地址"` // 文件地址
	Tag          string `json:"tag" gorm:"comment:文件标签"`                    // 文件标签
	Key          string `json:"key" gorm:"comment:编号"`                      // 编号
	Size         string `json:"size" gorm:"comment:size"`                   // 文件大小
	Reserve      string `json:"reserve" gorm:"comment:reserve"`             // 保留字段
}

type DetectionFileBatch struct {
	GVA_MODEL
	Batchid        string `json:"batchid" gorm:"comment:batchid"`
	Own            string `json:"own" gorm:"comment:Own"`                         // 拥有者
	App            string `json:"app" gorm:"comment:app"`                         // 应用
	FilesCount     int    `json:"files_count" gorm:"comment:files_count"`         // 文件数量
	FilesSize      string `json:"files_size" gorm:"comment:files_size"`           // 文件总大小
	BackendProgram string `json:"backend_program" gorm:"comment:backend_program"` // 后端程序
	Progress       string `json:"progress" gorm:"comment:progress"`               // 进度
	Status         string `json:"status" gorm:"comment:status"`                   // 状态
	BatchRank      string `json:"batch_rank" gorm:"comment:batch_rank"`           // 批次等级
	OutLabel       string `json:"out_label" gorm:"comment:out_label"`             // 输出标签
	Reserve        string `json:"reserve" gorm:"comment:reserve"`                 // 保留字段
}

func (DetectionFileUploadAndDownload) TableName() string {
	return "detection_file_upload_and_downloads"
}

type DetectionFileResponse struct {
	File DetectionFileUploadAndDownload `json:"file"`
}

func (DetectionFileBatch) TableName() string {
	return "detection_batches"
}

type DetectionBatchResponse struct {
	File DetectionFileBatch `json:"file"`
}
