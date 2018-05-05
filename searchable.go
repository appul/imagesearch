package imagesearch

import (
	"image"
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
		return newSearchableExact(img)
	}

	return newSearchableTolerance(img, tolerance)
}
