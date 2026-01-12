package server

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"pixgrid/converter"
	"strconv"
	"sync"
	"time"
)

type Session struct {
	Image     image.Image
	CreatedAt time.Time
	LastUsed  time.Time
}

type Server struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func New() *Server {
	s := &Server{
		sessions: make(map[string]*Session),
	}
	go s.cleanupLoop()
	return s
}

func (s *Server) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, session := range s.sessions {
			if now.Sub(session.LastUsed) > 30*time.Minute {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20) // 32MB max

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read image: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Failed to decode image: "+err.Error(), http.StatusBadRequest)
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "Failed to generate session ID", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	s.sessions[sessionID] = &Session{
		Image:     img,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}
	s.mu.Unlock()

	// Encode original image as base64 for preview
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"sessionId": sessionID,
		"width":     img.Bounds().Dx(),
		"height":    img.Bounds().Dy(),
		"original":  "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"sessionId"`
		Size      int    `json:"size"`
		Scale     int    `json:"scale"`
		Colors    int    `json:"colors"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	session, exists := s.sessions[req.SessionID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Update last used time
	s.mu.Lock()
	session.LastUsed = time.Now()
	s.mu.Unlock()

	// Apply defaults
	if req.Size <= 0 {
		req.Size = 64
	}
	if req.Scale <= 0 {
		req.Scale = 8
	}

	// Convert the image
	result := ConvertImage(session.Image, req.Size, req.Scale, req.Colors)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, result); err != nil {
		http.Error(w, "Failed to encode result", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"image":  "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
		"width":  result.Bounds().Dx(),
		"height": result.Bounds().Dy(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"sessionId"`
		Size      int    `json:"size"`
		Scale     int    `json:"scale"`
		Colors    int    `json:"colors"`
	}

	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	session, exists := s.sessions[req.SessionID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Apply defaults
	if req.Size <= 0 {
		req.Size = 64
	}
	if req.Scale <= 0 {
		req.Scale = 8
	}

	// Convert the image
	result := ConvertImage(session.Image, req.Size, req.Scale, req.Colors)

	// Encode to PNG and send as file
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=pixelart.png")
	png.Encode(w, result)
}

// ConvertImage applies the pixel art conversion to an in-memory image
func ConvertImage(img image.Image, pixelSize, scale, colors int) image.Image {
	// Downscale
	smallImg := converter.Downscale(img, pixelSize)

	// Quantize colors if specified
	if colors > 0 {
		smallImg = converter.QuantizeColors(smallImg, colors)
	}

	// Upscale with nearest neighbor
	return converter.UpscaleNearestNeighbor(smallImg, scale)
}

func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", s.corsMiddleware(s.handleUpload))
	mux.HandleFunc("/api/convert", s.corsMiddleware(s.handleConvert))
	mux.HandleFunc("/api/download", s.corsMiddleware(s.handleDownload))
	return mux
}

func (s *Server) Start(port int) error {
	mux := s.SetupRoutes()
	addr := ":" + strconv.Itoa(port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}
