package main

import (
	"errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/nickalie/go-webpbin"
)

func EncodeWebP(filePath string) (*string, error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	ext := path.Ext(filePath)

	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		return nil, errors.New("not an image")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("failed to open file")

		return nil, err
	}

	fileName := strings.Replace(filePath, ext, ".webp", -1)

	output, err := os.Create(fileName)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("failed to create file")

		return nil, err
	}
	defer output.Close()

	var img image.Image

	if ext == ".jpg" || ext == ".jpeg" {
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Error().Err(err).Str("file", filePath).Msg("failed to decode jpeg file")

			return nil, err
		}
	}
	if ext == ".png" {
		img, err = png.Decode(file)
		if err != nil {
			log.Error().Err(err).Str("file", filePath).Msg("failed to decode png file")

			return nil, err
		}
	}

	if err := webpbin.Encode(output, img); err != nil {
		output.Close()
		log.Error().Err(err).Str("file", filePath).Msg("failed to encode webp file")
	}

	return &fileName, nil
}

func DecodeWebP(filePath string) {

}
