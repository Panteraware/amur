package main

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ServeFile(c echo.Context) error {
	fileId, _ := url.PathUnescape(c.Request().URL.String())

	u, _ := url.Parse(fileId)
	result := u.Path

	filePath := Config.PublicFolder

	filePath = path.Clean(filePath + result)

	mType, err := mimetype.DetectFile(filePath)

	if err != nil {
		log.Error().Err(err).Str("path", "").Msg("error detecting mimetype")

		return c.NoContent(500)
	}

	if strings.Contains(mType.String(), "image") {
		if !strings.HasSuffix(filePath, ".webp") {
			go func() {
				_, err := EncodeWebP(filePath)
				if err != nil {
					log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error encoding webp")
				}
			}()
		}

		go func() {
			_, err := ResizeImage(filePath)
			if err != nil {
				log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error resizing image")
			}
		}()
	} else if strings.Contains(mType.String(), "video") {
		if Config.ConvertHLS {
			go func() {
				err := ConvertToHLS(filePath)
				if err != nil {
					log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error converting to hls")
				}
			}()
		}

		if len(Config.VideoScale) > 0 {
			fullPath, _ := filepath.Abs(filePath)
			//output, err := exec.Command(fmt.Sprintf("ffprobe -v error -select_streams v -show_entries stream=height -of csv=p=0:s=x \"%s\"", fullPath)).Output()
			//if err != nil {
			//	log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error executing ffprobe")
			//}
			//
			//videoScale, err := strconv.Atoi(string(output))

			for _, scale := range GetVideoScales(1080) {
				go func() {
					err := ScaleVideo(fullPath, scale)
					if err != nil {
						log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error scaling video")
					}
				}()
			}
		}
	}

	log.Info().Str("mType", mType.String()).Str("filePath", filePath).Str("user-agent", c.Request().UserAgent()).Str("ip", c.RealIP()).Msg("mimetype")

	buf, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Str("path", "").Msg("error reading file")

		return c.NoContent(500)
	}

	return c.Blob(200, mType.String(), buf)
}

func UploadFile(c echo.Context) error {
	key := c.Request().Header.Get("Authorization")

	if len(Config.UploadKey) < 8 {
		return c.NoContent(403)
	}

	if key != Config.UploadKey {
		return c.NoContent(400)
	}

	file, _ := c.FormFile("file")

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("error generating id")

		return c.NoContent(500)
	}

	ext := path.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", id, ext)
	filePath := path.Clean(fmt.Sprintf("%s/files/%s%s", Config.PublicFolder, id, ext))

	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("route", "upload").Msg("error opening file")
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(filePath)
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Str("route", "upload").Msg("error creating file")
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Error().Err(err).Str("path", filePath).Str("route", "upload").Msg("error copying file")
		return err
	}

	thumbUrl := ""

	if CheckFileExtension(fileName) == "image" {
		_, err := ResizeImage(filePath)

		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("error resizing file")

			return c.JSON(200, H{"url": fmt.Sprintf("%s/f/%s", Config.Domain, fileName)})
		}

		thumbUrl = fmt.Sprintf("%s_thumb%s", strings.Replace(fileName, ext, "", -1), ext)
	}

	return c.JSON(200, H{"url": fmt.Sprintf("%s/f/%s", Config.Domain, fileName), "thumb": fmt.Sprintf("%s/f/%s", Config.Domain, thumbUrl)})
}
