package ocr

import "net/http"

// OCRHttpHandler handles ocr http calls
type OCRHttpHandler struct{}

// ServeHTTP implements http.Handler interface
func (h *OCRHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
