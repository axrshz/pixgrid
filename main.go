package main

import (
	"flag"
	"fmt"
	"os"
	"pixie/converter"
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

	config := converter.Config{
		InputFile:  *inputFile,
		OutputFile: *outputFile,
		PixelSize:  *pixelSize,
		Scale:      *scale,
		Colors:     *colors,
	}

	if err := converter.Convert(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Conversion completed successfully!")
}