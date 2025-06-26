package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
)

func ServeFile(c echo.Context) error {
	fileId, _ := url.PathUnescape(c.Request().URL.String())

	u, _ := url.Parse(fileId)
	result := u.Path

	filePath := Config.PublicFolder

	filePath = path.Clean(filePath + result)

	exists := Exists(filePath)

	if !exists {
		return c.NoContent(404)
	}

	return c.File(filePath)
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
		_, err := ThumbnailImage(filePath)

		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("error resizing file")

			return c.JSON(200, H{"url": fmt.Sprintf("%s/f/%s", Config.Domain, fileName)})
		}

		thumbUrl = fmt.Sprintf("%s_thumb%s", strings.Replace(fileName, ext, "", -1), ext)
	}

	return c.JSON(200, H{"url": fmt.Sprintf("%s/f/%s", Config.Domain, fileName), "thumb": fmt.Sprintf("%s/f/%s", Config.Domain, thumbUrl)})
}
