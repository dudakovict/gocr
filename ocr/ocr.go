package ocr

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/dudakovict/gocr/conf"
	"github.com/otiai10/gosseract/v2"
)

// GosseractClient is an interface matching the methods used by gosseract.Client
type GosseractClient interface {
	SetImageFromBytes(content []byte) error
	Text() (string, error)
	Close() error
}

// OCRProcessor holds the OCR client and other related functionality
type OCRProcessor struct {
	gosseractClient GosseractClient
	logger          *log.Logger
}

// NewOCRProcessor initializes a new OCRProcessor
func NewOCRProcessor(logger *log.Logger) *OCRProcessor {
	return &OCRProcessor{
		gosseractClient: gosseract.NewClient(),
		logger:          logger,
	}
}

// UploadHandler handles the file upload and text extraction
func (ocr *OCRProcessor) UploadHandler(w http.ResponseWriter, r *http.Request, cfg conf.Config) {
	r.ParseMultipartForm(int64(cfg.MaxFileSizeMB) << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		ocr.logger.Printf("Error retrieving the file: %s", err)
		return
	}
	defer file.Close()

	if !isImage(handler) {
		http.Error(w, "Invalid file format. Only image files are allowed.", http.StatusBadRequest)
		ocr.logger.Println("Invalid file format. Only image files are allowed.")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		ocr.logger.Printf("Error reading the file: %s", err)
		return
	}

	if err := ocr.gosseractClient.SetImageFromBytes(fileBytes); err != nil {
		http.Error(w, "Error setting image from bytes", http.StatusInternalServerError)
		ocr.logger.Printf("Error setting image from bytes: %s", err)
		return
	}

	text, err := ocr.gosseractClient.Text()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error extracting text: %s", err), http.StatusInternalServerError)
		ocr.logger.Printf("Error extracting text: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, text)
}

// Close closes any resources held by the OCRProcessor
func (ocr *OCRProcessor) Close() error {
	return ocr.gosseractClient.Close()
}

// isImage checks if the given file is an image
func isImage(fileHeader *multipart.FileHeader) bool {
	if fileHeader == nil {
		return false
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		return false
	}

	return strings.HasPrefix(strings.ToLower(contentType), "image/")
}
