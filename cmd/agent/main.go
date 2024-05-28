package main

import (
	"github.com/iselldonuts/metrics/internal/config/agent"
	"github.com/iselldonuts/metrics/internal/core"
	"github.com/iselldonuts/metrics/internal/metrics"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer func() {
		_ = logger.Sync()
	}()
	log := logger.Sugar()

	cfg, err := getConfig()
	if err != nil {
		log.Panic(err)
	}
	run(cfg, log)
}

func run(conf *agent.Config, log *zap.SugaredLogger) {
	log.Infof(
		"Running agent | url: %s, ReportInterval: %d, PollInterval: %d\n",
		conf.Address, conf.ReportInterval, conf.PollInterval,
	)

	poller := metrics.NewPoller()
	a := core.NewAgent(poller, conf)
	a.Start(log)
}
