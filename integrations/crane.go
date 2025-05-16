package integrations

import "github.com/Tiktai/handler/model"

func CraneRetrieveLabels(elem model.ImageMetric) model.Image {
	return model.Image{Sha: elem.Sha, Name: "img-" + elem.Sha}
}
