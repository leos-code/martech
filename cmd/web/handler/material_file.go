package handler

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/ginx"
	"github.com/tencentad/martech/pkg/common/ffmpeg"
	"github.com/tencentad/martech/pkg/common/s3x"
	"github.com/gin-gonic/gin"
)

var materialFileOption *config.MaterialFileOption

// MaterialFileUploadHandler 素材文件上传接口
func MaterialFileUploadHandler(c *gin.Context) {
	if materialFileOption == nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("material_file option not configure"))
		return

	}
	var fileHeader *multipart.FileHeader
	var err error
	fileHeader, err = c.FormFile("file")
	if err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}
	size := fileHeader.Size

	var file multipart.File
	file, err = fileHeader.Open()
	if err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	filename := fileHeader.Filename
	uniqFileName := uniqFileName(filename)
	ext := getExt(filename)

	// 上传s3
	var url string
	url, err = s3x.GetS3().Upload(file, materialFileOption.Bucket, formatS3Key(uniqFileName), ext)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}
	var data *types.MaterialData
	data, err = processMaterial(&processMaterialParam{
		uniqFilename: uniqFileName,
		reader:       file,
		ext:          ext,
		size:         size,
		url:          url,
	})
	if err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	ginx.ResponseWithStatus200(c, "upload success", data)
}

func getExt(fullName string) string {
	pointIdx := strings.LastIndex(fullName, ".")
	if pointIdx == -1 {
		return ""
	}
	ext := fullName[pointIdx:]
	return strings.ToLower(ext)
}

func formatS3Key(filename string) string {
	return path.Join(
		materialFileOption.Path,
		filename)
}

func uniqFileName(filename string) string {
	return time.Now().Format("20060102150405") + "_" + filename
}

var (
	imageExt = map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
	}
)

func isImageExt(ext string) bool {
	_, ok := imageExt[ext]
	return ok
}

var (
	videoExt = map[string]struct{}{
		".mp4": {},
	}
)

func isVideoExt(ext string) bool {
	_, ok := videoExt[ext]
	return ok
}

type processMaterialParam struct {
	uniqFilename string
	reader       io.Reader
	ext          string
	size         int64
	url          string
}

func processMaterial(param *processMaterialParam) (*types.MaterialData, error) {
	if isImageExt(param.ext) {
		return processImageMaterial(param)
	} else if isVideoExt(param.ext) {
		return processVideoMaterial(param)
	} else {
		return nil, fmt.Errorf("not support ext: %s", param.ext)
	}
}

func processImageMaterial(param *processMaterialParam) (*types.MaterialData, error) {
	img, _, err := image.Decode(param.reader)
	if err != nil {
		return nil, err
	}

	return &types.MaterialData{
		Type: types.MaterialTypeImage,
		Image: &types.MaterialImage{
			Width:  img.Bounds().Size().X,
			Height: img.Bounds().Size().Y,
			Ext:    param.ext,
			URL:    param.url,
			Size:   param.size,
		},
	}, nil
}

func processVideoMaterial(param *processMaterialParam) (*types.MaterialData, error) {
	f, err := os.Create(param.uniqFilename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, param.reader)
	if err != nil {
		return nil, err
	}

	if err = f.Close(); err != nil {
		return nil, err
	}

	probeOutput, err := ffmpeg.ProbeWithTimeout(param.uniqFilename, 0)
	if err != nil {
		return nil, err
	}

	materialData := &types.MaterialData{
		Type: types.MaterialTypeVideo,
		Video: &types.MaterialVideo{
			URL:             param.url,
			Size:            param.size,
			ContainerFormat: probeOutput.Format.FormatName,
		},
	}

	av := materialData.Video
	for _, stream := range probeOutput.Streams {
		if stream.CodecType == ffmpeg.CodecTypeVideo {
			av.Duration = int(stream.Duration)
			av.CodecFormat = stream.CodecName
			av.Width = stream.Width
			av.Height = stream.Height
			av.Bitrate = int(stream.Bitrate / 1000)
			av.FrameRate = int(stream.FrameRate)
		} else if stream.CodecType == ffmpeg.CodecTypeAudio {
			av.AudioCodec = stream.CodecName
			av.AudioBitrate = int(stream.Bitrate / 1000)
			av.AudioSampleRate = float32(stream.SampleRate / 1000)
		}
	}

	return materialData, nil
}
