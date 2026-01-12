# pixgrid

Convert PNG/JPG images to pixel art.

## Install

```bash
go build
```

## CLI Usage

```bash
./pixgrid -input photo.jpg -output art.png
```

### Options

```
-input     Input image (required)
-output    Output file (default: output.png)
-size      Pixel width (default: 64)
-scale     Upscale factor (default: 8)
-colors    Color palette size, 0 to disable (default: 32)
```

### Examples

```bash
# Basic conversion
./pixgrid -input photo.jpg -output pixel.png

# Recommended
./pixgrid -input photo.jpg -output pixelart.png -size 64 -scale 8 -colors 32
```

## Web Interface

Pixgrid includes a web UI with real-time preview.

### Running the Web App

**1. Start the backend server:**

```bash
go run cmd/server/main.go
```

The server runs on `http://localhost:8080` by default. Use `-port` to change:

```bash
go run cmd/server/main.go -port 3000
```

**2. Start the frontend dev server:**

```bash
cd web
npm install
npm run dev
```

The frontend runs on `http://localhost:5173` and proxies API calls to the backend.

**3. Open `http://localhost:5173` in your browser**

### Features

- Drag-and-drop image upload
- Real-time preview as you adjust parameters
- Side-by-side comparison (original vs pixel art)
- One-click PNG download

### Building for Production

```bash
# Build the frontend
cd web
npm run build

# The build output is in web/dist/
```
