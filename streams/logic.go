package streams

import (
	"fmt"

	"github.com/Tiktai/handler/model"
)

// imageMetricToImage maps an ImageMetric to an Image for use in the stream

func imageMetricToImage(elem model.ImageMetric) model.Image {
	return model.Image{Sha: elem.Sha, Name: "img-" + elem.Sha}
}

type ImageWithError struct {
	Image model.Image
	Err   error
}

func imageMetricToImageWithError(elem model.ImageMetric) ImageWithError {
	if elem.Sha == "" {
		return ImageWithError{Err: fmt.Errorf("ImageMetric.Sha is empty")}
	}
	return ImageWithError{Image: model.Image{Sha: elem.Sha, Name: "img-" + elem.Sha}}
}
