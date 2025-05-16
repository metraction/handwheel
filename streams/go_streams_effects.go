package streams

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/reugn/go-streams/flow"
)

var counter int = 0

// fetchImageDetails maps an ImageMetric to an Image for use in the stream
//func fetchImageDetails(elem ImageMetric) chan struct{ Image } {

func fetchImageDetails(elem ImageMetric) Image {
	log.Println("start fetchImageDetails", elem.Sha, counter)
	seed := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(seed)
	duration := time.Duration(gen.Intn(5)) * time.Second

	result := Image{Sha: elem.Sha, Name: "img-" + strconv.Itoa(counter)}
	counter++
	time.Sleep(duration)
	log.Println("end fetchImageDetails", elem.Sha, " after ", duration)
	return result
}

// AugmentImages creates a stream that maps ImageMetric to Image using FlowMap and Via
func AugmentImages(metrics []ImageMetric) []Image {

	mapper := flow.NewMap(
		fetchImageDetails,
		1, // parallelism
	)

	go func() {
		for _, m := range metrics {
			mapper.In() <- m
		}
		close(mapper.In())
	}()

	outCh := mapper.Out()
	var images []Image
	for channel := range outCh {
		log.Println("Got channel reponse")
		//		out := <-channel.(chan struct{ Image })
		//		log.Println("Got AugmentImages reponse", out.Image)
		//		images = append(images, out.Image)
		images = append(images, channel.(Image))
	}
	return images
}
