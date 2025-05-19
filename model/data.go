package model

type ImageMetric struct {
	Image_spec string
}

type Image struct {
	Image_spec string
	Labels     map[string]string
}

type OutputWithError struct {
	Result any
	Err    error
}
