package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"
)

// TransformInput represents the input used to transform (compress, resize) a video
type TransformInput struct {
	Width          int64
	Height         int64
	QualityLevel   int64
	VideoCodec     string
	BitRate        int64
	InputFilename  string
	OutputFilename string
}

// VideoCodecs currently supported video codecs
var VideoCodecs = map[string]string{
	"h264": "libx264",
	"h265": "libx265",
}

// Transform applies the specified transformations on the video file
func Transform(input *TransformInput) error {
	if input.Height == 0 {
		input.Height = -1
	}

	if input.Width == 0 {
		input.Width = -1
	}

	commandLine := fmt.Sprintf(
		"ffmpeg -y -v error -hide_banner -i %s -vf scale=%d:%d",
		input.InputFilename,
		input.Width,
		input.Height,
	)
	commandLineParts := strings.Split(commandLine, " ")

	if input.QualityLevel > 1 {
		commandLineParts = append(commandLineParts, "-crf", fmt.Sprintf("%d", input.QualityLevel))
	}

	if input.VideoCodec != "" {
		lib, exists := VideoCodecs[input.VideoCodec]
		if !exists {
			return fmt.Errorf("invalid codec (%s) specified", input.VideoCodec)
		}

		commandLineParts = append(commandLineParts, "-vcodec", lib)

	}

	if input.BitRate > 0 {
		commandLineParts = append(commandLineParts, "-b", fmt.Sprintf("%dk", input.BitRate/1000))
	}

	commandLineParts = append(commandLineParts, input.OutputFilename)

	output, err := exec.Command(commandLineParts[0], commandLineParts[1:]...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s : %s", err.Error(), string(output))
	}

	return nil
}
