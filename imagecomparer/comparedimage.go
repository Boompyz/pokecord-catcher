package imagecomparer

import (
	"image"
	"io"
	"math"

	_ "image/jpeg" // for jpeg files
	_ "image/png"  // for png files
)

const (
	thumbnailSize = 8
)

// ComparedImage contains information about an image
type ComparedImage struct {
	thumbnail [thumbnailSize][thumbnailSize][3]uint32
}

// NewComparedImage creates image from io.Reader and extracts info
func NewComparedImage(imageSource io.Reader) (*ComparedImage, error) {

	var thumbnail [thumbnailSize][thumbnailSize][3]uint32
	var div [thumbnailSize][thumbnailSize]uint32

	m, _, err := image.Decode(imageSource)
	if err != nil {
		return nil, err
	}

	bounds := m.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA() // each in the range [0, 65535]

			tx := (x - bounds.Min.X) * thumbnailSize / bounds.Dx()
			ty := (y - bounds.Min.Y) * thumbnailSize / bounds.Dy()

			thumbnail[ty][tx][0] += r
			thumbnail[ty][tx][1] += g
			thumbnail[ty][tx][2] += b

			div[ty][tx]++
		}
	}

	for y := 0; y < thumbnailSize; y++ {
		for x := 0; x < thumbnailSize; x++ {
			for i := 0; i < 3; i++ {
				thumbnail[y][x][i] /= div[y][x]
			}
		}
	}

	return &ComparedImage{thumbnail}, nil
}

// GetDistance returns the difference from another image
func (c *ComparedImage) GetDistance(o *ComparedImage) float64 {
	if c == nil || o == nil {
		return 2147483647
	}
	var distance float64
	for y := 0; y < thumbnailSize; y++ {
		for x := 0; x < thumbnailSize; x++ {
			dr := float64(0)
			for i := 0; i < 3; i++ {
				d := c.thumbnail[y][x][i] - o.thumbnail[y][x][i]
				dr += float64(d * d)
			}
			distance += math.Sqrt(dr)
		}
	}

	return distance
}
