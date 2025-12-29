package converter

import (
	"image"
	"image/color"
)

func QuantizeColors(img image.Image, numColors int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	levelsPerChannel := int(float64(numColors) / 3.0)
	if levelsPerChannel < 2 {
		levelsPerChannel = 2
	}

	step := 255 / (levelsPerChannel - 1)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldColor := img.At(x, y)
			r, g, b, a := oldColor.RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			r8 = quantizeChannel(r8, step)
			g8 = quantizeChannel(g8, step)
			b8 = quantizeChannel(b8, step)

			newColor := color.RGBA{R: r8, G: g8, B: b8, A: a8}
			newImg.Set(x, y, newColor)
		}
	}

	return newImg
}

func quantizeChannel(value uint8, step int) uint8 {
	level := int(float64(value)/float64(step) + 0.5)
	result := level * step

	if result > 255 {
		result = 255
	}

	return uint8(result)
}