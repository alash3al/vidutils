package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alash3al/vidutils/pkg/ffmpeg"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "vidutils",
		Usage:                "a very simple video utilities utilizing the power of ffmpeg",
		EnableBashCompletion: true,
	}

	// inspect
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "inspect",
		Usage: "fetch some information about the specified media file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src",
				Usage:    "specifies the source file path",
				Aliases:  []string{"s"},
				Required: true,
			},
			&cli.BoolFlag{
				Name: "thumbnail",
			},
			&cli.StringFlag{
				Name:  "thumbnail-time-offset",
				Value: "00:00:01",
			},
			&cli.Int64Flag{
				Name:  "thumbnail-width",
				Value: 480,
			},
			&cli.Int64Flag{
				Name:  "thumbnail-height",
				Value: 360,
			},
		},
		Action: func(ctx *cli.Context) error {
			info, err := ffmpeg.Inspect(&ffmpeg.InspectInput{
				Filename:            ctx.String("src"),
				ExtractThumbnail:    ctx.Bool("thumbnail"),
				ThumbnailTimeOffset: ctx.String("thumbnail-time-offset"),
				ThumbnailWidth:      ctx.Int64("thumbnail-width"),
				ThumbnailHeight:     ctx.Int64("thumbnail-height"),
			})

			if err != nil {
				return err
			}

			fmt.Println(info.String())

			return nil
		},
	})

	// transform
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "transform",
		Usage: "scale, convert or change the media file quality",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src",
				Usage:    "specifies the source file path",
				Aliases:  []string{"s"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "out",
				Usage:    "specifies the output file path",
				Aliases:  []string{"o"},
				Required: true,
			},
			&cli.Int64Flag{
				Name:  "width",
				Usage: "the new output width",
				Value: -1,
			},
			&cli.Int64Flag{
				Name:  "height",
				Usage: "the new output height",
				Value: -1,
			},
			&cli.StringFlag{
				Name:    "video-codec",
				Aliases: []string{"vc"},
				Usage:   "the new output video-codec, currently we support (h264, h265)",
				Value:   "h264",
			},
			&cli.Int64Flag{
				Name:    "quality-level",
				Aliases: []string{"q", "crf"},
				Usage:   "the new output quality level (CRF 'Constant Rate Factor') ranged from 0 to 52 the lower value the more quality",
				Value:   -1,
			},
			&cli.Int64Flag{
				Name:    "video-bitrate",
				Aliases: []string{"vb"},
				Usage:   "the new output bitrate",
				Value:   -1,
			},
			&cli.Int64Flag{
				Name:    "audio-bitrate",
				Aliases: []string{"ab"},
				Usage:   "the new output bitrate",
				Value:   -1,
			},
			&cli.Int64Flag{
				Name:  "fps",
				Usage: "frames per second",
				Value: 30,
			},
		},
		Action: func(ctx *cli.Context) error {
			return ffmpeg.Transform(&ffmpeg.TransformInput{
				InputFilename:  ctx.String("src"),
				OutputFilename: ctx.String("out"),
				OutputQuality: ffmpeg.OutputQuality{
					Width:            ctx.Int64("width"),
					Height:           ctx.Int64("height"),
					QualityLevel:     ctx.Int64("quality-level"),
					VideoBitRateKilo: ctx.Int64("video-bitrate"),
					AudioBitRateKilo: ctx.Int64("audio-bitrate"),
				},
				VideoCodec: ctx.String("video-codec"),
				FrameRate:  ctx.Int64("fps"),
			})
		},
	})

	// hlsify
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "hlsify",
		Usage: "generate a hls playlist for the specified video",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src",
				Usage:    "the source video file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "out",
				Usage:    "the output directory",
				Required: true,
			},
			&cli.Int64Flag{
				Name:    "segment-duration",
				Aliases: []string{"d"},
				Usage:   "the segment duration in seconds",
				Value:   2,
			},
			&cli.StringSliceFlag{
				Name:     "quality-presets",
				Aliases:  []string{"q"},
				Usage:    "the required quality preset, available presets are: " + strings.Join(ffmpeg.QualityPresetsList, ","),
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			playlistPath, err := ffmpeg.GenerateHLSPlaylist(&ffmpeg.HLSBuilderInput{
				InputFilename:          ctx.String("src"),
				OutputDirectory:        ctx.String("out"),
				SegmentDurationSeconds: ctx.Int64("segment-duration"),
				QualityPresets:         ctx.StringSlice("quality-presets"),
			})

			fmt.Println(playlistPath)

			return err
		},
	})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
