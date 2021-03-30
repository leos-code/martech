package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*
{
    "streams": [
        {
            "index": 0,
            "codec_name": "h264",
            "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
            "profile": "Main",
            "codec_type": "video",
            "codec_time_base": "1/50",
            "codec_tag_string": "avc1",
            "codec_tag": "0x31637661",
            "width": 1280,
            "height": 720,
            "coded_width": 1280,
            "coded_height": 720,
            "closed_captions": 0,
            "has_b_frames": 0,
            "sample_aspect_ratio": "1:1",
            "display_aspect_ratio": "16:9",
            "pix_fmt": "yuv420p",
            "level": 31,
            "chroma_location": "left",
            "refs": 1,
            "is_avc": "true",
            "nal_length_size": "4",
            "r_frame_rate": "25/1",
            "avg_frame_rate": "25/1",
            "time_base": "1/12800",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 67584,
            "duration": "5.280000",
            "bit_rate": "1205959",
            "bits_per_raw_sample": "8",
            "nb_frames": "132",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "creation_time": "1970-01-01T00:00:00.000000Z",
                "language": "und",
                "handler_name": "VideoHandler",
                "vendor_id": "[0][0][0][0]"
            }
        },
        {
            "index": 1,
            "codec_name": "aac",
            "codec_long_name": "AAC (Advanced Audio Coding)",
            "profile": "LC",
            "codec_type": "audio",
            "codec_time_base": "1/48000",
            "codec_tag_string": "mp4a",
            "codec_tag": "0x6134706d",
            "sample_fmt": "fltp",
            "sample_rate": "48000",
            "channels": 6,
            "channel_layout": "5.1",
            "bits_per_sample": 0,
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/48000",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 254976,
            "duration": "5.312000",
            "bit_rate": "384828",
            "max_bit_rate": "400392",
            "nb_frames": "249",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "creation_time": "1970-01-01T00:00:00.000000Z",
                "language": "und",
                "handler_name": "SoundHandler",
                "vendor_id": "[0][0][0][0]"
            }
        }
    ],
    "format": {
        "filename": "sample.mp4",
        "nb_streams": 2,
        "nb_programs": 0,
        "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
        "format_long_name": "QuickTime / MOV",
        "start_time": "0.000000",
        "duration": "5.312000",
        "size": "1055736",
        "bit_rate": "1589963",
        "probe_score": 100,
        "tags": {
            "major_brand": "isom",
            "minor_version": "512",
            "compatible_brands": "isomiso2avc1mp41",
            "creation_time": "1970-01-01T00:00:00.000000Z",
            "encoder": "Lavf53.24.2"
        }
    }
}
*/

// ProbeOutput ffprobe 输出的信息
type ProbeOutput struct {
	Format  *Format   `json:"format"`
	Streams []*Stream `json:"streams"`
}

type CodecType string

const (
	CodecTypeVideo CodecType = "video"
	CodecTypeAudio CodecType = "audio"
)

type Stream struct {
	CodecType  CodecType    `json:"codec_type"`
	CodecName  string       `json:"codec_name"`
	Width      int          `json:"width"`
	Height     int          `json:"height"`
	Duration   StringNumber `json:"duration"`
	Bitrate    StringNumber `json:"bit_rate"`
	FrameRate  FrameRate    `json:"r_frame_rate"`
	SampleRate StringNumber `json:"sample_rate"`
}

type FrameRate int

func (f *FrameRate) UnmarshalJSON(d []byte) error {
	s := string(d)
	s = strings.Trim(s, "\"")
	splits := strings.Split(s, "/")
	if len(splits) == 2 {
		v, err := strconv.Atoi(splits[0])
		if err != nil {
			return err
		}

		*f = FrameRate(v)
	} else {
		*f = 0
	}
	return nil
}

type Format struct {
	Size       StringNumber `json:"size"`
	FormatName string       `json:"format_name"`
}

// StringNumber 字符串表示的数字，如"123"
type StringNumber float64

func (f *StringNumber) UnmarshalJSON(d []byte) error {
	s := string(d)
	s = strings.Trim(s, "\"")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*f = StringNumber(v)
	return nil
}

// ProbeWithTimeout 获取视频文件的信息
func ProbeWithTimeout(fileName string, timeOut time.Duration) (*ProbeOutput, error) {
	args := []string{"-show_format", "-show_streams", "-of", "json"}
	args = append(args, fileName)
	ctx := context.Background()
	if timeOut > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), timeOut)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, "ffprobe", args...)
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	output := &ProbeOutput{}
	if err = json.Unmarshal(buf.Bytes(), output); err != nil {
		return nil, err
	}

	return output, nil
}
