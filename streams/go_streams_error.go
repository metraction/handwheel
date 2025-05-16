package streams

import (
	"fmt"
	"github.com/reugn/go-streams/flow"
	"github.com/Tiktai/handler/model"
)

func ErrorRouter(metrics model.ImageMetric) model.Image {

	mapper := flow.NewMap(
		imageMetricToImageWithError,
		1, // parallelism
	)

	go func() {
		mapper.In() <- metrics
		close(mapper.In())
	}()

	result := <-mapper.Out()
	imgWithErr, ok := result.(ImageWithError)
	if !ok {
		panic("unexpected type from stream")
	}
	if imgWithErr.Err != nil {
		fmt.Printf("Error in ImageMetricToImageWithError: %v\n", imgWithErr.Err)
		return model.Image{}
	}
	return imgWithErr.Image
}
