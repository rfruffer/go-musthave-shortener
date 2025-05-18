package handlers

import (
	"encoding/json"
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
	body, err := io.ReadAll(c.Request.Body)
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

func (us *URLHandler) GetShortURLHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.String(http.StatusBadRequest, "missing ID")
		return
	}

	originalURL, err := us.service.RedirectURL(id)
	if err != nil {
		c.String(http.StatusBadRequest, "cant find id in store")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}

func (us *URLHandler) SetResultHost(host string) {
	us.baseURL = host
}
