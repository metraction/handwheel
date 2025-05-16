package streams

import (
	"github.com/Tiktai/handler/model"
	"github.com/reugn/go-streams/flow"
)

// ProcessImage creates a stream that maps ImageMetric to Image using FlowMap and Via
func ProcessImage(metrics []model.ImageMetric) []model.Image {

	mapper := flow.NewMap(
		imageMetricToImage,
		1, // parallelism
	)

	go func() {
		for _, m := range metrics {
			mapper.In() <- m
		}
		close(mapper.In())
	}()

	outCh := mapper.Out()
	var images []model.Image
	for img := range outCh {
		images = append(images, img.(model.Image))
	}
	return images
}
