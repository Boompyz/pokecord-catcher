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
	Thumbnail [thumbnailSize][thumbnailSize][3]uint32 `json:"thumbnail"`
}

// NewComparedImage creates image from io.Reader and extracts info
func NewComparedImage(imageSource io.Reader) (*ComparedImage, error) {

	var thumbnail [thumbnailSize][thumbnailSize][3]uint32
	var div [thumbnailSize][thumbnailSize]uint32

	m, _, err := image.Decode(imageSource)
	if err != nil {
		return nil, err
	}

	m = cropImage(m)

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
				d := c.Thumbnail[y][x][i] - o.Thumbnail[y][x][i]
				dr += float64(d * d)
			}
			distance += math.Sqrt(dr)
		}
	}

	return distance
}

func cropImage(m image.Image) image.Image {
	var top, left, right, bottom int
	top, left, right, bottom = m.Bounds().Min.Y, m.Bounds().Min.X, m.Bounds().Max.X, m.Bounds().Max.Y

	scanHorizontal := func(lineNumber int) bool {
		var y = lineNumber
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 || a != 0 {
				return false
			}
		}
		return true
	}
	scanVertical := func(colNumber int) bool {
		var x = colNumber
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			r, g, b, a := m.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 || a != 0 {
				return false
			}
		}
		return true
	}

	for scanHorizontal(top) {
		top++
	}

	for scanHorizontal(bottom) {
		bottom--
	}

	for scanVertical(left) {
		left++
	}

	for scanVertical(right) {
		right--
	}

	newImage := image.NewRGBA(image.Rect(left, top, right, bottom))
	for y := top; y < bottom; y++ {
		for x := left; x < right; x++ {
			newImage.Set(x, y, m.At(x, y))
		}
	}

	return newImage
}
