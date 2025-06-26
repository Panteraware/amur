package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var watcher *fsnotify.Watcher

type FileHistory struct {
	FilePath  string
	Operation string
	Time      time.Time
}

func WatchFolder() {
	history := []FileHistory{}

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(Config.PublicFolder, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				history = append(history, FileHistory{FilePath: event.Name, Operation: event.Op.String(), Time: time.Now()})

				if !strings.HasSuffix(event.Name, ".part") && event.Op.String() == "CHMOD" && history[len(history)-1].Operation != "REMOVE" {
					filePath := path.Clean(event.Name)

					mType := http.DetectContentType(ReadLimitedBytes(filePath, 512))

					log.Info().Str("file", event.Name).Str("type", mType).Msg("file created or modified")

					if strings.Contains(mType, "image") {
						if Config.UseRedis {
							if !strings.HasSuffix(filePath, ".webp") {
								task, err := NewImageOptimizationTask(filePath)
								if err != nil {
									log.Error().Err(err).Str("path", filePath).Str("route", "serve").Str("task", TypeImageOptimization).Msg("error creating task")
								}
								info, err := AsynqClient.Enqueue(task)
								if err != nil {
									log.Error().Err(err).Str("path", filePath).Str("route", "serve").Str("enqueue", TypeImageOptimization).Msg("error enqueuing task")
								}

								log.Info().Str("id", info.ID).Str("task", TypeImageOptimization).Msg("queued task")
							}

							task, err := NewImageThumbnailTask(filePath)
							if err != nil {
								log.Error().Err(err).Str("path", filePath).Str("route", "serve").Str("task", TypeImageThumbnail).Msg("error creating task")
							}
							info, err := AsynqClient.Enqueue(task)
							if err != nil {
								log.Error().Err(err).Str("path", filePath).Str("route", "serve").Str("enqueue", TypeImageThumbnail).Msg("error enqueuing task")
							}

							log.Info().Str("id", info.ID).Str("task", TypeImageThumbnail).Msg("queued task")
						} else {
							if !strings.HasSuffix(filePath, ".webp") {
								_, err := EncodeWebP(filePath)
								if err != nil {
									log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error encoding webp")
								}
							}

							_, err := ThumbnailImage(filePath)
							if err != nil {
								log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error thumbnailing image")
							}
						}
					} else if strings.Contains(mType, "video") {
						if Config.CanConvertHLS {
							go func() {
								err := ConvertToHLS(filePath)
								if err != nil {
									log.Error().Err(err).Str("path", filePath).Str("route", "serve").Msg("error converting to hls")
								}
							}()
						}

						if Config.CanScaleVideo && len(Config.VideoScale) > 0 {
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
				}

				history = []FileHistory{}

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
