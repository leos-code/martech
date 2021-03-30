package s3x

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

var (
	once      sync.Once
	singleton *S3Wrapper
)

// S3Wrapper
type S3Wrapper struct {
	s3       *s3.S3
	uploader *s3manager.Uploader
}

// S3Type s3类型
type S3Type string

const (
	S3TypeCos S3Type = "cos"
	S3TypeS3  S3Type = "s3"
)

// Option s3/cos 初始化选项
type Option struct {
	Type   S3Type `json:"type"`
	Region string `json:"region"`
	AppID  string `json:"app_id"`
	AppKey string `json:"app_key"`
}

// GetS3 获取S3单例
func GetS3(options ...*Option) *S3Wrapper {
	once.Do(func() {
		var option *Option
		if len(options) != 0 {
			option = options[0]
		}

		if option == nil {
			log.Errorf("call GetS3 without option for the first time")
			return
		}

		var err error
		singleton, err = create(option)
		if err != nil {
			log.Errorf("failed to create S3Wrapper, err: %v", err)
		}
	})

	return singleton
}

func create(option *Option) (*S3Wrapper, error) {
	var sess *session.Session

	if option.Type == S3TypeS3 {
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(option.Region),
		}))
	} else if option.Type == S3TypeCos {
		cred := credentials.NewStaticCredentials(option.AppID, option.AppKey, "")
		sess = session.Must(session.NewSession(&aws.Config{
			Credentials:      cred,
			Region:           aws.String(option.Region),
			EndpointResolver: endpoints.ResolverFunc(cosResolver),
		}))

	} else {
		log.Errorf("not valid s3")
		return nil, fmt.Errorf("not valid s3 type[%s]", option.Type)
	}

	return &S3Wrapper{
		s3:       s3.New(sess),
		uploader: s3manager.NewUploader(sess),
	}, nil
}

func cosResolver(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	if service == "s3" {
		return endpoints.ResolvedEndpoint{
			URL:           "http://cos." + region + ".myqcloud.com",
			SigningRegion: region,
		}, nil
	}
	return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
}

func (w *S3Wrapper) PutObject(bucket string, localPath, key string) error {
	fp, err := os.Open(localPath)
	if err != nil {
		log.Error("can not open file ", localPath)
		return err
	}
	defer fp.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()
	_, err = w.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   fp,
	})
	if err != nil {
		log.Error("put object error: ", err)
		return err
	}

	return nil
}

func (w *S3Wrapper) GetObject(bucket string, key string) (*s3.GetObjectOutput, error) {
	return w.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

// HeadObject 获取文件头信息
func (w *S3Wrapper) HeadObject(bucket string, key string) (*s3.HeadObjectOutput, error) {
	return w.s3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

// Upload 上传文件
func (w *S3Wrapper) Upload(file io.Reader, bucket string, key string, ext string) (string, error) {
	auth := "public-read"

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "binary/octet-stream"
	}
	result, err := w.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		ContentType: &mimeType,
		Key:         aws.String(key),
		Body:        file,
		ACL:         &auth,
	})

	if err != nil {
		log.Error("upload object error: ", err)
		return "", fmt.Errorf("failed to upload %s to %s, err: %v", key, bucket, err)
	}

	return result.Location, nil
}

// GetS3FilePathPart
func GetS3FilePathPart(path string) (string, string, string) {
	prefix := "://"
	index := strings.Index(path, prefix)
	if index == -1 {
		return "", "", ""
	}
	noPrefixPath := path[index+len(prefix):]

	index1 := strings.Index(noPrefixPath, "/")
	index2 := strings.LastIndex(noPrefixPath, "/")
	if index1 == -1 || index2 == -1 {
		return "", "", ""
	}
	var bucket, folder, filename string
	bucket = noPrefixPath[:index1]
	filename = noPrefixPath[index2+1:]

	if index2 > index1 {
		folder = noPrefixPath[index1+1 : index2]
	}

	return bucket, folder, filename
}
