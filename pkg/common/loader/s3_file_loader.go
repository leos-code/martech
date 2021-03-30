package loader

import (
	"path"

	"github.com/tencentad/martech/pkg/common/s3x"
	log "github.com/sirupsen/logrus"
)

// S3FileLoader s3/cos文件加载器
type S3FileLoader struct {
	bucket   string
	key      string
	alloc    allocFunc
	loadFunc fileLoadFunc

	lastModify int64
}

// DetectNewFile implement FileLoader
func (loader *S3FileLoader) DetectNewFile() (string, bool) {
	modify, err := loader.getModifyTime()
	if err != nil {
		return "", false
	}

	return "", modify > loader.lastModify
}

func (loader *S3FileLoader) getModifyTime() (int64, error) {
	head, err := s3x.GetS3().HeadObject(loader.bucket, loader.key)
	if err != nil {
		log.Warnf("failed to head object, err: %v", err)
		return 0, err
	}

	return head.LastModified.Unix(), nil
}

// Load implement FileLoader
func (loader *S3FileLoader) Load(_ string, i interface{}) error {
	modify, err := loader.getModifyTime()
	if err != nil {
		return err
	}

	output, err := s3x.GetS3().GetObject(loader.bucket, loader.key)
	if err != nil {
		return err
	}

	if err = loader.loadFunc(output.Body, i); err != nil {
		return err
	}

	loader.lastModify = modify
	return nil
}

// Alloc implement FileLoader
func (loader *S3FileLoader) Alloc() interface{} {
	return loader.alloc()
}

// NewS3FileLoader
func NewS3FileLoader(filepath string, alloc allocFunc, loadFunc fileLoadFunc) *S3FileLoader {
	bucket, folder, filename := s3x.GetS3FilePathPart(filepath)

	return &S3FileLoader{
		bucket:   bucket,
		key:      path.Join(folder, filename),
		alloc:    alloc,
		loadFunc: loadFunc,
	}
}
