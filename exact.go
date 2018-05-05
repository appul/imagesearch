package imagesearch

import "image"

// exactImage is a Searchable image that needs to match exactly
type exactImage struct {
	// width and height are the width and the height of the source
	width, height int

	// baseOffset is the offset of the first value in the source Pix slice.
	baseOffset int

	// img is the source img
	img *image.RGBA
}

func newExactImage(img image.Image) *exactImage {
	bounds := img.Bounds()

	// Make a copy of the source image
	src := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			src.Set(x, y, img.At(x, y))
		}
	}

	return &exactImage{
		img:        image.NewRGBA(bounds),
		width:      bounds.Dx(),
		height:     bounds.Dy(),
		baseOffset: src.PixOffset(bounds.Min.X, bounds.Min.Y),
	}
}

func (e *exactImage) SearchIn(img *image.RGBA) image.Rectangle {
	bounds := img.Bounds()

	// baseOffset is the offset of the first pixel of img in img.Pix
	baseOffset := img.PixOffset(bounds.Min.Y, bounds.Min.X)

	// The width and the height of the searchable region account for the width
	// of the source image
	width := bounds.Max.X - e.width + 1
	height := bounds.Max.Y - e.height + 1

	// Iterate over the pixels in the image
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {

			// Check if the region matches our source
			if !e.check(baseOffset+y*img.Stride+x*4, img) {
				continue
			}

			// Return a new rect from the current position with the width and
			// height of the source image
			return image.Rect(x, y, x+e.width, y+e.height)
		}
	}

	return image.ZR
}

// check compares a region with the source image and returns true if the pixels
// match exactly
func (e *exactImage) check(imgRegionOffset int, img *image.RGBA) bool {
	// Iterate over the pixels in the source and the img
	for y := 0; y < e.height; y++ {
		for x := 0; x < e.width; x++ {

			// Calculate the offsets for position (x,y) for the Pix slices
			imgOffset := imgRegionOffset + img.Stride*y + x*4
			srcOffset := e.baseOffset + e.img.Stride*y + x*4

			// Compare the pixel values of the image against the pixel values of
			// the source
			if img.Pix[imgOffset+0] != e.img.Pix[srcOffset+0] ||
				img.Pix[imgOffset+1] != e.img.Pix[srcOffset+1] ||
				img.Pix[imgOffset+2] != e.img.Pix[srcOffset+2] {
				return false
			}
		}
	}

	return true
}
