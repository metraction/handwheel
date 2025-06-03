package routing

import (
	"log"
	"time"

	"github.com/metraction/handwheel/integrations"
	"github.com/metraction/handwheel/model"
	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

func NewTicker(period time.Duration) chan any {
	outChan := make(chan any)
	ticker := time.NewTicker(period)
	go func() {
		for {
			outChan <- ""
			<-ticker.C
		}
	}()
	return outChan
}

func ProtheusCraneDevLakeRouter(cfg *model.Config) {

	interval, err := time.ParseDuration(cfg.Prometheus.Interval)
	if err != nil {
		log.Printf("Invalid prometheus.interval in config: %v, defaulting to 1m", err)
		interval = time.Minute
	}
	log.Printf("Using interval: %v", interval)
	source := ext.NewChanSource(NewTicker(interval))
	sink := ext.NewStdoutSink()

	prometheusIntegration := integrations.NewPrometheusIntegration(cfg)
	craneIntegration := integrations.NewCraneIntegration(cfg)
	devlakeIntegration := integrations.NewDevLakeIntegration(cfg)
	seenImages := integrations.NewDedup()
	/*
		TODO: detect image change which happened when application was not working
	*/
	source.
		Via(flow.NewMap(prometheusIntegration.FetchImageMetrics, 1)).
		Via(flow.NewFlatMap(integrations.PrometheusResult, 1)).
		Via(flow.NewFilter(craneIntegration.WhiteListImages(), 1)).
		Via(flow.NewFilter(seenImages.FilterDublicates, 1)).
		Via(flow.NewMap(craneIntegration.CraneRetrieveLabels, 1)).
		Via(flow.NewMap(devlakeIntegration.PostDeployment, 1)).
		To(sink)
}
