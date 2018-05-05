package imagesearch

import (
	"image"
	"image/color"
)

// toleranceImage is a Searchable image with a tolerance level
type toleranceImage struct {
	// width and height are the width and the height of the source
	width, height int

	// minPix and maxPix contain the values for the min and max tolerable
	// colors
	minPix, maxPix []uint8

	// baseOffset is the offset of the first value in the source Pix slices.
	baseOffset int

	// stride is the stride used for getting the offsets for minPix and maxPix
	stride int
}

func newToleranceImage(img image.Image, t uint8) *toleranceImage {
	bounds := img.Bounds()

	// Create images with the minimum and maximum tolerable color values for
	// comparison.
	minSrc, maxSrc := image.NewRGBA(bounds), image.NewRGBA(bounds)

	// Iterate over the pixels of the original source image and populate the min
	// and max images
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Combine the transparency (255 - alpha) with the tolerance to get
			// the actual tolerance
			tolerance := toleranceMax(0xFF-uint8(a), t)

			// Set the min color
			minSrc.SetRGBA(x, y, color.RGBA{
				R: toleranceMin(uint8(r), tolerance),
				G: toleranceMin(uint8(g), tolerance),
				B: toleranceMin(uint8(b), tolerance),
			})

			// Set the max color
			maxSrc.SetRGBA(x, y, color.RGBA{
				R: toleranceMax(uint8(r), tolerance),
				G: toleranceMax(uint8(g), tolerance),
				B: toleranceMax(uint8(b), tolerance),
			})
		}
	}

	// Precalculate variables that will remain consistent to improve performance
	return &toleranceImage{
		width:      bounds.Dx(),
		height:     bounds.Dy(),
		minPix:     minSrc.Pix,
		maxPix:     maxSrc.Pix,
		baseOffset: minSrc.PixOffset(bounds.Min.X, bounds.Min.Y),
		stride:     minSrc.Stride,
	}
}

func (t *toleranceImage) SearchIn(img *image.RGBA) image.Rectangle {
	bounds := img.Bounds()

	// baseOffset is the offset of the first pixel of img in img.Pix
	baseOffset := img.PixOffset(bounds.Min.Y, bounds.Min.X)

	// The width and the height of the searchable region account for the width
	// of the source image
	width := bounds.Max.X - t.width + 1
	height := bounds.Max.Y - t.height + 1

	// Iterate over the pixels in the image
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {

			// Check if the region matches our source
			if !t.check(baseOffset+y*img.Stride+x*4, img) {
				continue
			}

			// Return a new rect from the current position with the width and
			// height of the source image
			return image.Rect(x, y, x+t.width, y+t.height)
		}
	}

	return image.ZR
}

// check compares a region with the source image and returns true if it's a
// match within the tolerable range
func (t *toleranceImage) check(imgRegionOffset int, img *image.RGBA) bool {
	// Iterate over the pixels in the source and the img
	for y := 0; y < t.height; y++ {
		for x := 0; x < t.width; x++ {

			// Calculate the offsets for position (x,y) for the Pix slices
			imgOffset := imgRegionOffset + img.Stride*y + x*4
			srcOffset := t.baseOffset + t.stride*y + x*4

			// Compare the pixel values against the min and max tolerable pixel
			// values
			if img.Pix[imgOffset+0] < t.minPix[srcOffset+0] ||
				img.Pix[imgOffset+1] < t.minPix[srcOffset+1] ||
				img.Pix[imgOffset+2] < t.minPix[srcOffset+2] ||
				img.Pix[imgOffset+0] > t.maxPix[srcOffset+0] ||
				img.Pix[imgOffset+1] > t.maxPix[srcOffset+1] ||
				img.Pix[imgOffset+2] > t.maxPix[srcOffset+2] {
				return false
			}
		}
	}

	return true
}

// toleranceMax returns the max value that is within the tolerance of a given
// value.
func toleranceMax(v, tolerance uint8) uint8 {
	if tolerance > 0xFF-v {
		return 0xFF
	}
	return v + tolerance
}

// toleranceMin returns the min value that is within the tolerance of a given
// value.
func toleranceMin(v, tolerance uint8) uint8 {
	if tolerance > v {
		return 0
	}
	return v - tolerance
}
