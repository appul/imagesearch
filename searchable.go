package imagesearch

import (
	"image"
	"image/png"
	"os"
)

// A Searchable can be searched for in an image
type Searchable interface {
	// SearchIn searches for the Searchable in an image and return the Rectangle
	// that covers the Searchable
	SearchIn(*image.RGBA) image.Rectangle
}

// NewSearchableImage returns a new Searchable using an image to search for
func NewSearchableImage(img image.Image, tolerance uint8) Searchable {
	if tolerance == 0 {
		return newExactImage(img)
	}

	return newToleranceImage(img, tolerance)
}

// LoadSearchablePng loads an image from a PNG file and returns a new Searchable
// using that image
func LoadSearchablePng(file string, tolerance uint8) (Searchable, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	return NewSearchableImage(img, tolerance), nil
}
