package streams

import (
	"strconv"
	"testing"
)

func TestProcessImage(t *testing.T) {
	// Generate batches of ImageMetrics
	metrics := []ImageMetric{}
	for i := 0; i < 10; i++ {
		metrics = append(metrics, ImageMetric{Sha: strconv.Itoa(i)})
	}

	images := ProcessImage(metrics)

	if len(images) != len(metrics) {
		t.Fatalf("Expected %d images, got %d", len(metrics), len(images))
	}
	for i, img := range images {
		expectedSha := strconv.Itoa(i)
		expectedName := "img-" + expectedSha
		if img.Sha != expectedSha || img.Name != expectedName {
			t.Errorf("Unexpected mapping at index %d: got (Sha=%s, Name=%s), want (Sha=%s, Name=%s)", i, img.Sha, img.Name, expectedSha, expectedName)
		}
	}
}
