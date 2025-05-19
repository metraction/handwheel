package routing

import (
	"log"
	"time"

	"github.com/Tiktai/handler/integrations"
	"github.com/Tiktai/handler/model"
	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

func NewTicker(period time.Duration) chan any {
	outChan := make(chan any)
	ticker := time.NewTicker(period)
	go func() {
		for {
			<-ticker.C
			outChan <- ""
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
	craneIntegration := integrations.NewCraneIntegration(cfg.Crane)

	source.
		Via(flow.NewMap(prometheusIntegration.FetchImageMetrics, 1)).
		Via(flow.NewFlatMap(integrations.PrometheusResult, 1)).
		Via(flow.NewFilter(craneIntegration.WhiteListImages(), 1)).
		Via(flow.NewMap(craneIntegration.CraneRetrieveLabels, 1)).
		To(sink)
}
