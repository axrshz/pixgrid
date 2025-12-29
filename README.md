# pixgrid

convert png/jpg images to pixel art.

## Install

```bash
go build
```

## Usage

```bash
./pixgrid -input photo.jpg -output art.png
```

## Options

```
-input     Input image (required)
-output    Output file (default: output.png)
-size      Pixel width (default: 64)
-scale     Upscale factor (default: 8)
-colors    Color palette size, 0 to disable (default: 32)
```

## Examples

```bash
# Basic conversion
./pixgrid -input photo.jpg -output pixel.png

# Recommended
./pixgrid -input photo.jpg -output pixelart.png -size 64 -scale 8 -colors 32
```
