package main

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"time"
)

const (
	TypeImageOptimization = "image:optimization"
	TypeImageThumbnail    = "image:thumbnail"
	TypeImageResize       = "image:resize"
	TypeVideoTranscode    = "video:transcode"
	TypeVideoResize       = "video:resize"
)

type ImageOptimizationPayload struct {
	FilePath string
}

type ImageThumbnailPayload struct {
	FilePath string
}

type ImageResizePayload struct {
	FilePath string
	Width    int
	Height   int
}

type VideoTranscodePayload struct{}

type VideoResizePayload struct{}

func NewImageOptimizationTask(filePath string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{FilePath: filePath})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageOptimization, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

func NewImageThumbnailTask(filePath string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageThumbnailPayload{FilePath: filePath})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageThumbnail, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

func NewImageResizeTask(filePath string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{FilePath: filePath})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageResize, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

//func NewVideoTranscodeTask(filePath string) (*asynq.Task, error) {
//	payload, err := json.Marshal(VideoTranscodePayload{FilePath: filePath})
//	if err != nil {
//		return nil, err
//	}
//	// task options can be passed to NewTask, which can be overridden at enqueue time.
//	return asynq.NewTask(TypeImageOptimization, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
//}
//
//func NewVideoResizeTask() (*asynq.Task, error) {
//
//}
