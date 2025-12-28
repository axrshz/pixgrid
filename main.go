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

	// Save the image (no processing yet, just testing our pipeline)
	err = saveImage(*outputFile, img)
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