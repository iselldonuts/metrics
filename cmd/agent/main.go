package main

import (
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/iselldonuts/metrics/internal/config/agent"
	"github.com/iselldonuts/metrics/internal/core"
	"github.com/iselldonuts/metrics/internal/metrics"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		println(err)
		os.Exit(1)
	}
	log := logger.Sugar()

	conf, err := getConfig()
	if err != nil {
		_ = logger.Sync()
		log.Fatal(err)
	}
	log.Infow("Config loaded", "config", conf)

	run(conf, log)
}

func run(conf *agent.Config, log *zap.SugaredLogger) {
	log.Infof(
		"Running agent | url: %s, ReportInterval: %d, PollInterval: %d\n",
		conf.Address, conf.ReportInterval, conf.PollInterval,
	)

	p := metrics.NewPoller()
	s := metrics.NewSender(conf.Address, resty.New(), log)
	a := core.NewAgent(p, s, conf, log)
	a.Start()
}
