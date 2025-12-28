package main

import (
	"flag"
	"fmt"
	"image"
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

	// Save the downscaled image
	err = saveImage(*outputFile, smallImg)
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