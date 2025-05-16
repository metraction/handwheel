package routing

import (
	"github.com/Tiktai/handler/integrations"
	"github.com/Tiktai/handler/model"
	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

func ProtheusCraneDevLakeRouter(cfg *model.Config) {
	metricsChan := integrations.FetchAndSubmitPeriodically(cfg)
	source := ext.NewChanSource(metricsChan)
	mapFlow := flow.NewMap(integrations.CraneRetrieveLabels, 1)
	sink := ext.NewStdoutSink()

	source.
		Via(flow.Flatten[model.ImageMetric](1)).
		Via(mapFlow).
		To(sink)
}
