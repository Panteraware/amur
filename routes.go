package main

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"path"
	"strings"
)

func ServeFile(c *gin.Context) {
	fileId, _ := url.PathUnescape(c.Request.URL.String())

	u, _ := url.Parse(fileId)
	result := u.Path

	filePath := Config.PublicFolder

	filePath = path.Clean(filePath + result)

	mType, err := mimetype.DetectFile(filePath)

	if err != nil {
		log.Error().Err(err).Str("path", "").Msg("error detecting mimetype")

		c.Status(500)
		return
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
	}

	log.Info().Str("mType", mType.String()).Str("filePath", filePath).Str("user-agent", c.Request.UserAgent()).Str("ip", c.RemoteIP()).Msg("mimetype")

	buf, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Str("path", "").Msg("error reading file")

		c.Status(500)
		return
	}

	c.Data(200, mType.String(), buf)
}

func UploadFile(c *gin.Context) {
	key := c.GetHeader("Authorization")

	if len(Config.UploadKey) < 8 {
		c.Status(403)
		return
	}

	if key != Config.UploadKey {
		c.Status(400)
		return
	}

	file, _ := c.FormFile("file")

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("error generating id")

		c.Status(500)
		return
	}

	ext := path.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", id, ext)
	filePath := path.Clean(fmt.Sprintf("%s/files/%s%s", Config.PublicFolder, id, ext))

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("error saving file")

		c.Status(500)
		return
	}

	thumbUrl := ""

	if CheckFileExtension(fileName) == "image" {
		_, err := ResizeImage(filePath)

		if err != nil {
			log.Error().Err(err).Str("path", filePath).Str("util", "resize").Msg("error resizing file")

			c.JSON(200, gin.H{"url": fmt.Sprintf("%s/f/%s", Config.Domain, fileName)})
			return
		}

		thumbUrl = fmt.Sprintf("%s_thumb%s", strings.Replace(fileName, ext, "", -1), ext)
	}

	c.JSON(200, gin.H{"url": fmt.Sprintf("%s/f/%s", Config.Domain, fileName), "thumb": fmt.Sprintf("%s/f/%s", Config.Domain, thumbUrl)})
}
