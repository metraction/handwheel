package streams

import (
	"strconv"
	"testing"

	"github.com/Tiktai/handler/model"
)

func TestRequestReply(t *testing.T) {
	// Generate batches of ImageMetrics

	img := RequestReplyRouter(model.ImageMetric{Sha: strconv.Itoa(1)})
	expectedSha := strconv.Itoa(1)
	expectedName := "img-" + expectedSha
	if img.Sha != expectedSha || img.Name != expectedName {
		t.Errorf("Unexpected mapping at index %d: got (Sha=%s, Name=%s), want (Sha=%s, Name=%s)", 1, img.Sha, img.Name, expectedSha, expectedName)
	}
}
