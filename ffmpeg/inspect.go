package ffmpeg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/rs/xid"
	"github.com/tidwall/gjson"
)

// StreamMeta represents a stream in a media file
type StreamMeta struct {
	Index       int64     `json:"index"`
	Type        string    `json:"type"`
	Codec       string    `json:"codec"`
	Width       int64     `json:"width,omitempty"`
	Height      int64     `json:"height,omitempty"`
	AspectRatio string    `json:"aspect_ratio,omitempty"`
	StartTime   float64   `json:"start_time"`
	Duration    float64   `json:"duration"`
	BitRate     int64     `json:"bit_rate"`
	FramesCount int64     `json:"frames_count,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// InspectOutput represents a media file info
type InspectOutput struct {
	Duration float64 `json:"duration"`
	Size     struct {
		Bytes float64 `json:"bytes"`
		Human string  `json:"human"`
	} `json:"size"`
	BitRate   int64         `json:"bit_rate"`
	Thumbnail string        `json:"thumbnail,omitempty"`
	Streams   []*StreamMeta `json:"streams"`
}

func (info *InspectOutput) String() string {
	j, _ := json.Marshal(info)

	return string(j)
}

// InspectInput the input used to extract the media information
type InspectInput struct {
	Filename            string
	ExtractThumbnail    bool
	ThumbnailTimeOffset string
	ThumbnailWidth      int64
	ThumbnailHeight     int64
}

// Inspect returns a media file information
func Inspect(input *InspectInput) (*InspectOutput, error) {
	commandLine := fmt.Sprintf("ffprobe -v error -print_format json -show_error -show_format -show_streams %s", input.Filename)
	commandLineParts := strings.Split(commandLine, " ")

	output, err := exec.Command(commandLineParts[0], commandLineParts[1:]...).CombinedOutput()
	outputParsed := gjson.ParseBytes(output)

	if err != nil {
		if !outputParsed.Get("error.string").Exists() {
			return nil, err
		}
		return nil, errors.New(outputParsed.Get("error.string").String())
	}

	info := &InspectOutput{
		Duration: outputParsed.Get("format.duration").Float(),
		BitRate:  outputParsed.Get("format.bit_rate").Int(),
	}

	info.Size.Bytes = outputParsed.Get("format.size").Float()
	info.Size.Human = units.HumanSizeWithPrecision(float64(info.Size.Bytes), 2)

	outputParsed.Get("streams").ForEach(func(_, value gjson.Result) bool {
		meta := StreamMeta{
			Index:       value.Get("index").Int(),
			Codec:       value.Get("codec_name").String(),
			Type:        value.Get("codec_type").String(),
			Width:       value.Get("width").Int(),
			Height:      value.Get("height").Int(),
			AspectRatio: value.Get("display_aspect_ratio").String(),
			StartTime:   value.Get("start_time").Float(),
			Duration:    value.Get("duration").Float(),
			BitRate:     value.Get("bit_rate").Int(),
			FramesCount: value.Get("nb_frames").Int(),
		}

		if value.Get("tags.creation_time").Exists() {
			meta.CreatedAt = value.Get("tags.creation_time").Time()
		}

		info.Streams = append(info.Streams, &meta)

		return true
	})

	if input.ExtractThumbnail && input.ThumbnailTimeOffset != "" && input.ThumbnailWidth > 0 && input.ThumbnailHeight > 0 {
		thumbnailFilename := filepath.Join(os.TempDir(), fmt.Sprintf("thumbnail-%s.png", xid.New().String()))
		commandLine := fmt.Sprintf(
			"ffmpeg -v error -hide_banner -i %s -ss %s -s %dx%d -vframes 1 -c:v png %s",
			input.Filename,
			input.ThumbnailTimeOffset,
			input.ThumbnailWidth,
			input.ThumbnailHeight,
			thumbnailFilename,
		)
		commandLineParts := strings.Split(commandLine, " ")
		output, err := exec.Command(commandLineParts[0], commandLineParts[1:]...).CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("%s : %s", err.Error(), string(output))
		}
		info.Thumbnail = thumbnailFilename
	}

	return info, nil
}
