package main

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func ResizeImage(filePath string, width int, height int) (string, error) {
	input, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to open file")
		return "", err
	}
	defer input.Close()

	fileExt := filepath.Ext(filePath)

	output, err := os.Create(fmt.Sprintf("%s_thumb%s", strings.Replace(filePath, fileExt, "", -1), fileExt))
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to create file")
		return "", err
	}
	defer output.Close()

	var (
		src image.Image
	)

	if fileExt == ".png" {
		src, err = png.Decode(input)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to decode png")
			return "", err
		}
	} else if fileExt == ".jpg" || fileExt == ".jpeg" {
		src, err = jpeg.Decode(input)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to decode jpeg")
			return "", err
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	if fileExt == ".png" {
		err := png.Encode(output, dst)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to encode png")
			return "", err
		}
	} else if fileExt == ".jpg" || fileExt == ".jpeg" {
		err := jpeg.Encode(output, dst, nil)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to encode jpeg")
			return "", err
		}
	}

	return fmt.Sprintf("%s_thumb%s", strings.Replace(filePath, fileExt, "", -1), fileExt), nil
}

func ThumbnailImage(filePath string) (string, error) {
	if strings.Contains(filePath, "_thumb") {
		return "", nil
	}

	input, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to open file")
		return "", err
	}
	defer input.Close()

	fileExt := filepath.Ext(filePath)

	output, err := os.Create(fmt.Sprintf("%s_thumb%s", strings.Replace(filePath, fileExt, "", -1), fileExt))
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to create file")
		return "", err
	}
	defer output.Close()

	var (
		src image.Image
	)

	if fileExt == ".png" {
		src, err = png.Decode(input)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to decode png")
			return "", err
		}
	} else if fileExt == ".jpg" || fileExt == ".jpeg" {
		src, err = jpeg.Decode(input)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to decode jpeg")
			return "", err
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))

	draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	if fileExt == ".png" {
		err := png.Encode(output, dst)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to encode png")
			return "", err
		}
	} else if fileExt == ".jpg" || fileExt == ".jpeg" {
		err := jpeg.Encode(output, dst, nil)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to encode jpeg")
			return "", err
		}
	}

	return fmt.Sprintf("%s_thumb%s", strings.Replace(filePath, fileExt, "", -1), fileExt), nil
}

func CheckFileExtension(filePath string) string {
	if strings.HasSuffix(filePath, ".webp") {
		return "image"
	} else if strings.HasSuffix(filePath, ".png") {
		return "image"
	} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
		return "image"
	}

	return ""
}

func GetVideoScales(size int) []string {
	initialArray := strings.Split(Config.VideoScale, ",")
	endArray := make([]string, 0)

	for _, v := range initialArray {
		if v == "2160" && size >= 2160 {
			endArray = append(endArray, "3840:2160")
		} else if v == "1440" && size >= 1440 {
			endArray = append(endArray, "2560:1440")
		} else if v == "1080" && size >= 1080 {
			endArray = append(endArray, "1920:1080")
		} else if v == "720" && size >= 720 {
			endArray = append(endArray, "1280:720")
		} else if v == "540" && size >= 540 {
			endArray = append(endArray, "960:540")
		} else if v == "480" && size >= 480 {
			endArray = append(endArray, "854:480")
		} else if v == "360" && size >= 360 {
			endArray = append(endArray, "640:360")
		} else if v == "240" && size >= 240 {
			endArray = append(endArray, "426:240")
		}
	}

	return endArray
}

func ReadLimitedBytes(path string, n int) []byte {
	exists := Exists(path)
	if !exists {
		return make([]byte, 0)
	}

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes := make([]byte, n)
	m, err := file.Read(bytes)
	if err != nil {
		panic(err)
	}

	return bytes[:m]
}
