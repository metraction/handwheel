package model

// A metric as returned by Prometheus, Labels are namespace, cluster and so on.
type ImageMetric struct {
	Image_spec string
	Labels     map[string]string
}

// An image as returned by Crane - Labels are different then the ones in Prometheus
type Image struct {
	Image_spec string
	Labels     map[string]string
}

type OutputWithError struct {
	Result any
	Err    error
}
