package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"path/filepath"
	"strings"
)

func ConvertToHLS(filePath string) error {
	hlsFilePath := strings.Replace(filePath, filepath.Ext(filePath), ".m3u8", 1)

	err := ffmpeg.Input(filePath).
		Output(hlsFilePath, ffmpeg.KwArgs{
			"b:v":           "1M",
			"g":             60,
			"hls_time":      10,
			"hls_list_size": 0,
			//"hls_segment_size": 500000,
		}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()

	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("util", "hls").Msg("convert to hls failed")
		return err
	}

	return nil
}

func ScaleVideo(filePath string, scale string) error {
	shrinkFilePath := fmt.Sprintf("%s_%s_%s", strings.Replace(filePath, filepath.Ext(filePath), "", 1), strings.Replace(scale, ":", "x", -1), filepath.Ext(filePath))

	err := ffmpeg.Input(filePath).
		Output(shrinkFilePath, ffmpeg.KwArgs{
			"vf": "scale=" + scale,
		}).
		OverWriteOutput().
		//ErrorToStdOut().
		Run()

	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("scale", scale).Str("util", "scale").Msg("shrink video failed")
		return err
	}

	return nil
}
