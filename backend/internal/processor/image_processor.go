package processor

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"time"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) ProcessImage(url string) (float64, error) {
	// Download image
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	// Decode image to get dimensions
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %w", err)
	}

	// Calculate perimeter
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	perimeter := 2 * float64(width+height)

	// Random sleep to simulate GPU processing (0.1 to 0.4 seconds)
	sleepTime := 100 + rand.Intn(301) // 100-400 milliseconds
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return perimeter, nil
}
