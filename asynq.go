package main

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

var AsynqClient *asynq.Client

func NewAsynqServer() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     Config.RedisHost,
			Password: Config.RedisPass,
		},
		asynq.Config{
			Concurrency: 100,
			Queues: map[string]int{
				"critical": 50,
				"default":  35,
				"low":      15,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeImageOptimization, HandleImageOptimization)
	mux.HandleFunc(TypeImageThumbnail, HandleImageThumbnail)
	mux.HandleFunc(TypeImageResize, HandleImageResize)

	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msg("failed to start asynq server")
	}
}

func NewAsynqClient() {
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     Config.RedisHost,
		Password: Config.RedisPass,
	})
}
