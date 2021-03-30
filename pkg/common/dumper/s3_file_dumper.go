package dumper

import (
	"path"

	"github.com/tencentad/martech/pkg/common/s3x"
	log "github.com/sirupsen/logrus"
)

type S3FileDumper struct {
	bucket   string
	folder   string
	filename string
}

func NewS3FileDumper(filepath string) FileDumper {
	bucket, folder, filename := s3x.GetS3FilePathPart(filepath)
	return &S3FileDumper{bucket: bucket, folder: folder, filename: filename}
}

func (d *S3FileDumper) Dump(obj interface{}, marshal Marshal) error {
	localFile := d.filename + ".tmp"
	if err := DumpToFile(obj, marshal, localFile); err != nil {
		return err
	}

	err := s3x.GetS3().PutObject(d.bucket, localFile, path.Join(d.folder, d.filename))
	if err != nil {
		log.Errorf("S3PutObject failed: %v", err)
		return err
	}
	return nil
}
