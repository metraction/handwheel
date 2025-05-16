package streams

import (
	"strconv"
	"testing"

	"github.com/Tiktai/handler/model"
)

func TestAugmentImage(t *testing.T) {
	// Generate batches of ImageMetrics
	metrics := []model.ImageMetric{}
	for i := 0; i < 3; i++ {
		metrics = append(metrics, model.ImageMetric{Sha: strconv.Itoa(i)})
	}

	images := AugmentImages(metrics)

	if len(images) != len(metrics) {
		t.Fatalf("Expected %d images, got %d", len(metrics), len(images))
	}
	for i, img := range images {
		expectedSha := strconv.Itoa(i)
		expectedName := "img-" + expectedSha
		t.Log(img)
		if img.Sha != expectedSha || img.Name != expectedName {
			t.Errorf("Unexpected mapping at index %d: got (Sha=%s, Name=%s), want (Sha=%s, Name=%s)", i, img.Sha, img.Name, expectedSha, expectedName)
		}
	}
}
