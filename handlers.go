package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func HandleImageOptimization(ctx context.Context, t *asynq.Task) error {
	var p ImageOptimizationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	_, err := EncodeWebP(p.FilePath)
	if err != nil {
		log.Error().Err(err).Str("path", p.FilePath).Str("route", "serve").Msg("error encoding webp")
		return err
	}

	return nil
}

func HandleImageResize(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	_, err := ResizeImage(p.FilePath, p.Width, p.Height)
	if err != nil {
		log.Error().Err(err).Str("path", p.FilePath).Str("route", "serve").Msg("error resizing image")
		return err
	}

	return nil
}

func HandleImageThumbnail(ctx context.Context, t *asynq.Task) error {
	var p ImageThumbnailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	_, err := ThumbnailImage(p.FilePath)
	if err != nil {
		log.Error().Err(err).Str("path", p.FilePath).Str("route", "serve").Msg("error resizing image")
		return err
	}

	return nil
}

func HandleVideoTranscode(ctx context.Context, t *asynq.Task) error {
	var p VideoTranscodePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}

func HandleVideoResize(ctx context.Context, t *asynq.Task) error {
	var p VideoResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}
