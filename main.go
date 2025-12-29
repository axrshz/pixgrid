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
	// Define command-line flags
	inputFile := flag.String("input", "", "Input image file (PNG or JPG)")
	outputFile := flag.String("output", "output.png", "Output image file")
	pixelSize := flag.Int("size", 64, "Target width in pixels (height scales proportionally)")
	scale := flag.Int("scale", 8, "Upscale factor (how much to enlarge the pixelated image)")
	colors := flag.Int("colors", 16, "Number of colors in the palette (0 = no quantization)")
	
	flag.Parse()

	// Validate input
	if *inputFile == "" {
		fmt.Println("Error: -input flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Load the image
	img, err := loadImage(*inputFile)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded image: %dx%d pixels\n", img.Bounds().Dx(), img.Bounds().Dy())

	// Step 1: Downscale the image to create the pixelated effect
	smallImg := downscale(img, *pixelSize)
	fmt.Printf("Downscaled to: %dx%d pixels\n", smallImg.Bounds().Dx(), smallImg.Bounds().Dy())

	// Step 2: Reduce colors if requested
	if *colors > 0 {
		smallImg = quantizeColors(smallImg, *colors)
		fmt.Printf("Reduced to %d colors\n", *colors)
	}

	// Step 3: Upscale back to a viewable size with hard pixel edges
	finalImg := upscaleNearestNeighbor(smallImg, *scale)
	fmt.Printf("Upscaled to: %dx%d pixels\n", finalImg.Bounds().Dx(), finalImg.Bounds().Dy())

	// Save the final image
	err = saveImage(*outputFile, finalImg)
	if err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Saved to: %s\n", *outputFile)
}

// loadImage reads an image file and decodes it
func loadImage(filename string) (image.Image, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Decode the image (automatically detects PNG or JPEG)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}

	return img, nil
}

// saveImage writes an image to a file
func saveImage(filename string, img image.Image) error {
	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	// Determine format from extension
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

// downscale reduces the image to a smaller size
func downscale(img image.Image, targetWidth int) image.Image {
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Calculate target height to maintain aspect ratio
	targetHeight := (origHeight * targetWidth) / origWidth

	// Create a new image with the target dimensions
	newImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Calculate the scaling factors
	scaleX := float64(origWidth) / float64(targetWidth)
	scaleY := float64(origHeight) / float64(targetHeight)

	// For each pixel in the new (small) image...
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// Find the corresponding pixel in the original image
			// We use the CENTER of the source pixel block for better quality
			srcX := int((float64(x) + 0.5) * scaleX)
			srcY := int((float64(y) + 0.5) * scaleY)

			// Get the color from the original image
			color := img.At(srcX, srcY)

			// Set it in the new image
			newImg.Set(x, y, color)
		}
	}

	return newImg
}

// quantizeColors reduces the number of colors in the image
func quantizeColors(img image.Image, numColors int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a new image for the quantized result
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Calculate how many levels per color channel
	// For example, 16 colors = ~2.5 levels per channel (2^4 = 16 total combinations)
	// We use cube root to distribute evenly across R, G, B
	levelsPerChannel := int(float64(numColors) / 3.0)
	if levelsPerChannel < 2 {
		levelsPerChannel = 2
	}

	// Calculate the step size for rounding
	// 255 is max color value, we divide into equal steps
	step := 255 / (levelsPerChannel - 1)

	// Process each pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the original color
			oldColor := img.At(x, y)
			r, g, b, a := oldColor.RGBA()

			// Convert from uint32 (0-65535) to uint8 (0-255)
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Round each color channel to the nearest step
			r8 = quantizeChannel(r8, step)
			g8 = quantizeChannel(g8, step)
			b8 = quantizeChannel(b8, step)

			// Create the new color and set it
			newColor := color.RGBA{R: r8, G: g8, B: b8, A: a8}
			newImg.Set(x, y, newColor)
		}
	}

	return newImg
}

// quantizeChannel rounds a color value to the nearest step
func quantizeChannel(value uint8, step int) uint8 {
	// Divide by step, round, then multiply back
	// Example: if step=51 and value=120
	// 120/51 = 2.35 → rounds to 2 → 2*51 = 102
	level := int(float64(value)/float64(step) + 0.5)
	result := level * step

	// Make sure we don't exceed 255
	if result > 255 {
		result = 255
	}

	return uint8(result)
}

// upscaleNearestNeighbor enlarges the image while keeping hard pixel edges
func upscaleNearestNeighbor(img image.Image, scaleFactor int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions
	newWidth := width * scaleFactor
	newHeight := height * scaleFactor

	// Create the larger image
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// For each pixel in the NEW large image...
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Figure out which pixel from the small image this should be
			// We divide by scaleFactor to map back to the small image
			srcX := x / scaleFactor
			srcY := y / scaleFactor

			// Get the color from the small image
			color := img.At(srcX, srcY)

			// Set it in the large image
			newImg.Set(x, y, color)
		}
	}

	return newImg
}