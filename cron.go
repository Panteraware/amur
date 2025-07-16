package main

import (
	"github.com/rs/zerolog/log"
	"path"
	"strings"

	"os"
	"path/filepath"
	"time"
)

func CronInit() {
	for {
		go CheckFiles()

		time.Sleep(5 * time.Minute)
	}
}

func CheckFiles() {
	err := filepath.Walk(Config.PublicFolder, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error().Err(err)
			return err
		}

		if !info.IsDir() {
			exists := Exists(strings.Replace(p, path.Ext(p), ".webp", -1))

			if !exists {
				file, err2 := EncodeWebP(p)
				if err2 != nil {
					log.Error().Err(err2)
					return nil
				}

				log.Info().Str("file", *file).Msg("encoded image")
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal().Err(err)
	}
}
