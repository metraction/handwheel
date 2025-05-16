package streams

import (
	"github.com/Tiktai/handler/model"
	"github.com/reugn/go-streams/flow"
)

func RequestReplyRouter(metrics model.ImageMetric) model.Image {

	mapper := flow.NewMap(
		imageMetricToImage,
		1, // parallelism
	)

	go func() {
		mapper.In() <- metrics
		close(mapper.In())
	}()

	img := <-mapper.Out()
	return img.(model.Image)
}
