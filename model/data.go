package model

type ImageMetric struct {
	Sha string
}

type Image struct {
	Sha  string
	Name string
}

type OutputWithError struct {
	Result any
	Err    error
}
