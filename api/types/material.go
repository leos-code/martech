package types

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

var (
	MaterialTemp = &Material{}
)

// Material 素材
type Material struct {
	ID           uint64              `gorm:"column:id;primaryKey"                     json:"id,omitempty"`
	CreatedAt    time.Time           `gorm:"column:created_at"                        json:"created_at,omitempty"`
	UpdatedAt    time.Time           `gorm:"column:updated_at"                        json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt      `gorm:"column:delete_at;index"                   json:"-"`
	Name         string              `gorm:"column:name"                              json:"name,omitempty"`
	Data         *MaterialData       `gorm:"column:data;"                             json:"data,omitempty"`
	Audit        []*MaterialAudit    `gorm:"foreignKey:MaterialID"                    json:"audit,omitempty"`
	AuditStatus  MaterialAuditStatus `gorm:"column:audit_status;default:unaudited"    json:"audit_status,omitempty"`
	RejectReason string              `gorm:"column:reject_reason"                     json:"reject_reason,omitempty"`
}

// MaterialType 素材类型
type MaterialType string

const (
	MaterialTypeImage MaterialType = "image"
	MaterialTypeVideo MaterialType = "video"
)

// MaterialData 素材数据
type MaterialData struct {
	Type  MaterialType   `json:"type,omitempty"`
	Image *MaterialImage `json:"image,omitempty"`
	Video *MaterialVideo `json:"video,omitempty"`
}

// Scan
func (d *MaterialData) Scan(value interface{}) error {
	return scan(value, d)
}

// Value
func (d *MaterialData) Value() (driver.Value, error) {
	return value(d)
}

// MaterialImage 图片素材
type MaterialImage struct {
	URL    string `json:"url"`
	Size   int64  `json:"size"`             // 图片大小，单位byte
	Ext    string `json:"ext"`              // 后缀格式
	Width  int    `json:"width,omitempty"`  // 图片宽，单位像素
	Height int    `json:"height,omitempty"` // 图片高，单位像素
}

// MaterialVideo 视频素材
type MaterialVideo struct {
	URL             string  `json:"url"`
	Size            int64   `json:"size"`
	Duration        int     `json:"duration"`                    // 视频时长，单位秒
	ContainerFormat string  `json:"container_format,omitempty"`  // 格式，如mp4
	CodecFormat     string  `json:"codec_format,omitempty"`      // 视频编码格式，如h.264
	Width           int     `json:"width,omitempty"`             // 视频分辨率，宽，单位像素
	Height          int     `json:"height,omitempty"`            // 视频分辨率，高，单位像素
	Bitrate         int     `json:"bitrate,omitempty"`           // 码率，单位kbit/s
	FrameRate       int     `json:"frame_rate,omitempty"`        // 帧率
	AudioCodec      string  `json:"audio_codec,omitempty"`       // 音频格式，如AAC
	AudioBitrate    int     `json:"audio_bitrate,omitempty"`     // 音频码率，单位kbit/s
	AudioSampleRate float32 `json:"audio_sample_rate,omitempty"` // 音频采样率，单位kHZ
}
