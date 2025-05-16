package streams

import (
	"github.com/reugn/go-streams/flow"
)

func RequestReplyRouter(metrics ImageMetric) Image {

	mapper := flow.NewMap(
		imageMetricToImage,
		1, // parallelism
	)

	go func() {
		mapper.In() <- metrics
		close(mapper.In())
	}()

	img := <-mapper.Out()
	return img.(Image)
}
