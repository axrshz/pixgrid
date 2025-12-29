package converter

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	InputFile  string
	OutputFile string
	PixelSize  int
	Scale      int
	Colors     int
}

func Convert(config Config) error {
	img, err := loadImage(config.InputFile)
	if err != nil {
		return fmt.Errorf("loading image: %w", err)
	}

	fmt.Printf("Loaded image: %dx%d pixels\n", img.Bounds().Dx(), img.Bounds().Dy())

	smallImg := Downscale(img, config.PixelSize)
	fmt.Printf("Downscaled to: %dx%d pixels\n", smallImg.Bounds().Dx(), smallImg.Bounds().Dy())

	if config.Colors > 0 {
		smallImg = QuantizeColors(smallImg, config.Colors)
		fmt.Printf("Reduced to %d colors\n", config.Colors)
	}

	finalImg := UpscaleNearestNeighbor(smallImg, config.Scale)
	fmt.Printf("Upscaled to: %dx%d pixels\n", finalImg.Bounds().Dx(), finalImg.Bounds().Dy())

	if err := saveImage(config.OutputFile, finalImg); err != nil {
		return fmt.Errorf("saving image: %w", err)
	}

	fmt.Printf("Saved to: %s\n", config.OutputFile)
	return nil
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