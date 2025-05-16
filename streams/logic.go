package streams

import "fmt"

// Input steam
type ImageMetric struct {
	Sha string
}

type Image struct {
	Sha  string
	Name string
}

// imageMetricToImage maps an ImageMetric to an Image for use in the stream
func imageMetricToImage(elem ImageMetric) Image {
	return Image{Sha: elem.Sha, Name: "img-" + elem.Sha}
}

type ImageWithError struct {
	Image Image
	Err   error
}

func imageMetricToImageWithError(elem ImageMetric) ImageWithError {
	if elem.Sha == "" {
		return ImageWithError{Err: fmt.Errorf("ImageMetric.Sha is empty")}
	}
	return ImageWithError{Image: Image{Sha: elem.Sha, Name: "img-" + elem.Sha}}
}
