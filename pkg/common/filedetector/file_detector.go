package filedetector

import "os"

// FileDetector 文件更新接口
type FileDetector interface {
	DetectNewFile() (string, bool, error)
}

// LocalFileDetector 本地文件更新检查
type LocalFileDetector struct {
	filePath       string
	lastModifyTime int64
}

// NewLocalFileDetector 创建本地文件更新检查器
func NewLocalFileDetector(filePath string) *LocalFileDetector {
	return &LocalFileDetector{
		filePath:       filePath,
		lastModifyTime: 0,
	}
}

// DetectNewFile 检查本地文件是否有更新
func (d *LocalFileDetector) DetectNewFile() (newPath string, new bool, err error) {
	var fileInfo os.FileInfo
	fileInfo, err = os.Stat(d.filePath)
	if err != nil {
		return
	}

	modTime := fileInfo.ModTime().Unix()
	if modTime > d.lastModifyTime {
		newPath = d.filePath
		new = true
	} else {
		new = false
	}

	return
}
