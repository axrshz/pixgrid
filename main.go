package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputFile := flag.String("input", "", "Input image file (PNG or JPG)")
	outputFile := flag.String("output", "output.png", "Output image file")
	pixelSize := flag.Int("size", 64, "Target width in pixels (height scales proportionally)")
	scale := flag.Int("scale", 8, "Upscale factor (how much to enlarge the pixelated image)")
	colors := flag.Int("colors", 16, "Number of colors in the palette (0 = no quantization)")
	
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: -input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	img, err := loadImage(*inputFile)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded image: %dx%d pixels\n", img.Bounds().Dx(), img.Bounds().Dy())

	smallImg := downscale(img, *pixelSize)
	fmt.Printf("Downscaled to: %dx%d pixels\n", smallImg.Bounds().Dx(), smallImg.Bounds().Dy())

	if *colors > 0 {
		smallImg = quantizeColors(smallImg, *colors)
		fmt.Printf("Reduced to %d colors\n", *colors)
	}

	finalImg := upscaleNearestNeighbor(smallImg, *scale)
	fmt.Printf("Upscaled to: %dx%d pixels\n", finalImg.Bounds().Dx(), finalImg.Bounds().Dy())

	err = saveImage(*outputFile, finalImg)
	if err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Saved to: %s\n", *outputFile)
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}

	return img, nil
}

func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".png":
		err = png.Encode(file, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	default:
		return fmt.Errorf("unsupported output format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("could not encode image: %w", err)
	}

	return nil
}

func downscale(img image.Image, targetWidth int) image.Image {
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

func quantizeColors(img image.Image, numColors int) image.Image {
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

func upscaleNearestNeighbor(img image.Image, scaleFactor int) image.Image {
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