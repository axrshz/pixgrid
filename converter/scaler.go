package converter

import "image"

func Downscale(img image.Image, targetWidth int) image.Image {
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	targetHeight := (origHeight * targetWidth) / origWidth

	newImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	scaleX := float64(origWidth) / float64(targetWidth)
	scaleY := float64(origHeight) / float64(targetHeight)

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			srcX := int((float64(x) + 0.5) * scaleX)
			srcY := int((float64(y) + 0.5) * scaleY)

			color := img.At(srcX, srcY)

			newImg.Set(x, y, color)
		}
	}

	return newImg
}

func UpscaleNearestNeighbor(img image.Image, scaleFactor int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newWidth := width * scaleFactor
	newHeight := height * scaleFactor

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x / scaleFactor
			srcY := y / scaleFactor

			color := img.At(srcX, srcY)

			newImg.Set(x, y, color)
		}
	}

	return newImg
}