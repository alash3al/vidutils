package ffmpeg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/valyala/fasttemplate"
)

// package wide var(s)
var (
	QualityPresetsList = []string{"144p", "240p", "360p", "480p", "720p", "720p+60fps", "1080p", "1080p+60fps", "4k", "4k+60fps"}
	QualityPresetsMap  = map[string]OutputQuality{
		"144p": {
			Width:            256,
			Height:           144,
			VideoBitRateKilo: 90,
			AudioBitRateKilo: 32,
			QualityLevel:     24,
			FrameRate:        30,
		},
		"240p": {
			Width:            426,
			Height:           240,
			VideoBitRateKilo: 300,
			AudioBitRateKilo: 64,
			QualityLevel:     28,
			FrameRate:        30,
		},
		"360p": {
			Width:            640,
			Height:           360,
			VideoBitRateKilo: 700,
			AudioBitRateKilo: 96,
			QualityLevel:     24,
			FrameRate:        30,
		},
		"480p": {
			Width:            850,
			Height:           480,
			VideoBitRateKilo: 1400,
			AudioBitRateKilo: 128,
			QualityLevel:     24,
			FrameRate:        30,
		},
		"720p": {
			Width:            1280,
			Height:           720,
			VideoBitRateKilo: 2850,
			AudioBitRateKilo: 128,
			QualityLevel:     24,
			FrameRate:        30,
		},
		"720p+60fps": {
			Width:            1280,
			Height:           720,
			VideoBitRateKilo: 3950,
			AudioBitRateKilo: 128,
			QualityLevel:     20,
			FrameRate:        60,
		},
		"1080p": {
			Width:            1920,
			Height:           1080,
			VideoBitRateKilo: 4900,
			AudioBitRateKilo: 192,
			QualityLevel:     20,
			FrameRate:        30,
		},
		"1080p+60fps": {
			Width:            1920,
			Height:           1080,
			VideoBitRateKilo: 660,
			AudioBitRateKilo: 192,
			QualityLevel:     20,
			FrameRate:        60,
		},
		"4k": {
			Width:            3840,
			Height:           2160,
			VideoBitRateKilo: 14000,
			AudioBitRateKilo: 192,
			QualityLevel:     18,
			FrameRate:        30,
		},
		"4k+60fps": {
			Width:            3840,
			Height:           2160,
			VideoBitRateKilo: 25000,
			AudioBitRateKilo: 192,
			QualityLevel:     18,
			FrameRate:        60,
		},
	}
)

// HLSBuilderInput represents the input for the HLS builder
type HLSBuilderInput struct {
	InputFilename          string
	OutputDirectory        string
	SegmentDurationSeconds int64
	Quality                map[string]OutputQuality
	QualityPresets         []string
}

// OutputQuality represents the quality configs
type OutputQuality struct {
	Width            int64
	Height           int64
	VideoBitRateKilo int64
	AudioBitRateKilo int64
	QualityLevel     int64
	FrameRate        int64
}

// GenerateHLSPlaylist generate a HLS playlist from a video input
func GenerateHLSPlaylist(input *HLSBuilderInput) (string, error) {
	if err := os.MkdirAll(input.OutputDirectory, 0775); err != nil && err != os.ErrExist {
		return "", err
	}

	playlistFilename := path.Join(input.OutputDirectory, "playlist.m3u8")
	playlistContent := []string{
		"#EXTM3U",
		"#EXT-X-VERSION:3",
	}

	tpl := fasttemplate.New(
		"-vf scale=w={{width}}:h={{height}}:force_original_aspect_ratio=decrease,fps={{fps}} -c:a aac -ar 48000 -c:v h264 -profile:v main -crf {{quality_level}} -sc_threshold 0 -g 48 -keyint_min 48 -hls_time {{segment_duration_seconds}} -hls_playlist_type vod -b:v {{video_bitrate}}k -maxrate {{maxrate}}k -bufsize {{buffer_size}}k -b:a {{audio_bitrate}}k -hls_segment_filename {{output_dir}}/{{name}}_%09d.ts {{output_dir}}/{{name}}.m3u8",
		"{{", "}}",
	)

	args := strings.Split("-max_error_rate 0.0 -y -v error -hide_banner -i", " ")
	args = append(args, input.InputFilename)

	if input.Quality == nil {
		input.Quality = map[string]OutputQuality{}
	}

	for _, name := range input.QualityPresets {
		preset, exists := QualityPresetsMap[name]
		if !exists {
			return "", fmt.Errorf("quality preset %s is undefined", name)
		}

		input.Quality[name] = preset
	}

	usedQualities := 0

	for _, name := range QualityPresetsList {
		if _, exists := QualityPresetsMap[name]; !exists {
			continue
		}

		if _, exists := input.Quality[name]; !exists {
			continue
		}

		usedQualities++

		quality := QualityPresetsMap[name]

		parsedArgs := tpl.ExecuteString(map[string]interface{}{
			"width":                    strconv.FormatInt(quality.Width, 10),
			"height":                   strconv.FormatInt(quality.Height, 10),
			"quality_level":            strconv.FormatInt(quality.QualityLevel, 10),
			"segment_duration_seconds": strconv.FormatInt(input.SegmentDurationSeconds, 10),
			"video_bitrate":            strconv.FormatInt(quality.VideoBitRateKilo, 10),
			"audio_bitrate":            strconv.FormatInt(quality.AudioBitRateKilo, 10),
			"buffer_size":              strconv.FormatInt(int64(float64(quality.VideoBitRateKilo)/0.66666), 10),
			"maxrate":                  strconv.FormatInt(int64(float64(quality.VideoBitRateKilo)/0.934579), 10),
			"output_dir":               input.OutputDirectory,
			"fps":                      strconv.FormatInt(quality.FrameRate, 10),
			"name":                     name,
		})

		args = append(args, strings.Split(parsedArgs, " ")...)
		playlistContent = append(playlistContent, fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d", quality.VideoBitRateKilo*100, quality.Width, quality.Height))
		playlistContent = append(playlistContent, fmt.Sprintf("%s.m3u8", name))
	}

	if usedQualities < 1 {
		return "", fmt.Errorf("please specify at least one valid quality/preset")
	}

	output, err := exec.Command("ffmpeg", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s : %s", err.Error(), string(output))
	}

	if err := ioutil.WriteFile(playlistFilename, []byte(strings.Join(playlistContent, "\n")), 0775); err != nil {
		return "", err
	}

	return playlistFilename, nil
}
