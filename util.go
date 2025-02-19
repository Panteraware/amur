package main

import (
	"fmt"
	"github.com/nickalie/go-webpbin"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func ResizeImage(filePath string) (string, error) {
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
	} else if fileExt == ".webp" {
		src, err = webp.Decode(input)
		if err != nil {
			log.Error().Err(err).Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to decode webp")
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
	} else if fileExt == ".webp" {
		err := webpbin.Encode(output, dst)
		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("failed to encode webp")
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
