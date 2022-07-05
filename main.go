package main

import (
	"time"

	"github.com/Xenfo/watcher/internal/config"
	"github.com/Xenfo/watcher/internal/core"
	"github.com/Xenfo/watcher/internal/logger"

	"github.com/go-co-op/gocron"
)

func main() {
	config, err := config.Create()
	if err != nil {
		panic(err)
	}

	logger, err := logger.Create()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	scheduler := core.CreateScheduler(config, logger, gocron.NewScheduler(time.UTC))

	scheduler.Start()
	defer scheduler.Stop()
}
