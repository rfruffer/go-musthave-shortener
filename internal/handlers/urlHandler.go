package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

type URLHandler struct {
	service *services.URLService
	baseURL string
}

func NewURLHandler(service *services.URLService, baseURL string) *URLHandler {
	return &URLHandler{service: service, baseURL: baseURL}
}

func (us *URLHandler) CreateShortJSONURLHandler(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "empty or invalid body")
		return
	}
	id, err := us.service.GenerateShortURL(req.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to create a short url")
		return
	}

	resp := models.ShortenResponse{
		Result: us.baseURL + "/" + id,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(c.Writer).Encode(resp)
	// c.JSON(http.StatusCreated, resp)
}

func (us *URLHandler) CreateShortURLHandler(c *gin.Context) {
	var reader io.Reader = c.Request.Body
	if c.Request.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "failed to read gzip body")
			return
		}
		defer gzReader.Close()
		reader = gzReader
	}

	body, err := io.ReadAll(reader)
	if err != nil || len(body) == 0 {
		c.String(http.StatusBadRequest, "empty or invalid body")
		return
	}
	originalURL := string(body)
	id, err := us.service.GenerateShortURL(originalURL)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to create a short url")
		return
	}

	shortURL := us.baseURL + "/" + id
	c.Data(http.StatusCreated, "text/plain", []byte(shortURL))
}

// func (us *URLHandler) CreateShortURLHandler(c *gin.Context) {
// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil || len(body) == 0 {
// 		c.String(http.StatusBadRequest, "empty or invalid body")
// 		return
// 	}
// 	originalURL := string(body)
// 	id, err := us.service.GenerateShortURL(originalURL)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "failed to create a short url")
// 		return
// 	}

// 	shortURL := us.baseURL + "/" + id
// 	c.Data(http.StatusCreated, "text/plain", []byte(shortURL))
// }

func (us *URLHandler) GetShortURLHandler(c *gin.Context) {
	id := c.Param("id")

	fmt.Println("=== DEBUG START ===")
	fmt.Println("[DEBUG] Method:", c.Request.Method)
	fmt.Println("[DEBUG] Host:", c.Request.Host)
	fmt.Println("[DEBUG] FullPath:", c.FullPath())
	fmt.Println("[DEBUG] RequestURI:", c.Request.RequestURI)
	fmt.Println("[DEBUG] Path:", c.Request.URL.Path)
	fmt.Println("[DEBUG] Param ID:", id)

	// Query-параметры
	fmt.Println("[DEBUG] Query Params:")
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Заголовки
	fmt.Println("[DEBUG] Headers:")
	for name, values := range c.Request.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}

	// Попробуем прочитать тело запроса (если вдруг кто-то туда что-то шлёт)
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("[DEBUG] Error reading body:", err)
	} else if len(bodyBytes) > 0 {
		fmt.Println("[DEBUG] Body:", string(bodyBytes))
	} else {
		fmt.Println("[DEBUG] Body: <empty>")
	}

	// Вернём тело обратно, чтобы middleware или gin могли его использовать дальше
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if id == "" {
		fmt.Println("[DEBUG] ID is empty")
		c.String(http.StatusBadRequest, "missing ID")
		fmt.Println("=== DEBUG END ===")
		return
	}

	originalURL, err := us.service.RedirectURL(id)
	if err != nil {
		fmt.Println("[DEBUG] RedirectURL error:", err)
		c.String(http.StatusBadRequest, "cant find id in store")
		fmt.Println("=== DEBUG END ===")
		return
	}

	fmt.Println("[DEBUG] Redirecting to:", originalURL)
	fmt.Println("=== DEBUG END ===")

	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}

// func (us *URLHandler) GetShortURLHandler(c *gin.Context) {
// 	id := c.Param("id")

// 	if id == "" {
// 		c.String(http.StatusBadRequest, "missing ID")
// 		return
// 	}

// 	originalURL, err := us.service.RedirectURL(id)
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "cant find id in store")
// 		return
// 	}
// 	c.Redirect(http.StatusTemporaryRedirect, originalURL)
// }

func (us *URLHandler) SetResultHost(host string) {
	us.baseURL = host
}
