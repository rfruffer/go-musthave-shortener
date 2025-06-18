package handlers

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rfruffer/go-musthave-shortener/internal/async"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
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
	userID := c.GetString("user_id")
	id, err := us.service.GenerateShortURL(req.URL, userID)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			resp := models.ShortenResponse{
				Result: us.baseURL + "/" + id,
			}
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.WriteHeader(http.StatusConflict)
			json.NewEncoder(c.Writer).Encode(resp)
			return
		}
		c.String(http.StatusInternalServerError, "failed to create a short url")
		return
	}

	resp := models.ShortenResponse{
		Result: us.baseURL + "/" + id,
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(c.Writer).Encode(resp)
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
	userID := c.GetString("user_id")
	id, err := us.service.GenerateShortURL(originalURL, userID)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.Data(http.StatusConflict, "text/plain", []byte(us.baseURL+"/"+id))
			return
		}
		c.String(http.StatusInternalServerError, "failed to create a short url")
		fmt.Println("originalURL " + originalURL)
		fmt.Println("userID " + userID)
		fmt.Println("GenerateShortURL " + id)
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
		if errors.Is(err, repository.ErrGone) {
			c.String(http.StatusGone, "url deleted")
			return
		}
		c.String(http.StatusBadRequest, "cant find id in store")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}

func (us *URLHandler) SetResultHost(host string) {
	us.baseURL = host
}

func (us *URLHandler) Ping(c *gin.Context) {
	if err := us.service.Ping(); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}

func (us *URLHandler) Batch(c *gin.Context) {
	var req []models.BatchOriginalURL
	if err := c.ShouldBindJSON(&req); err != nil || len(req) == 0 {
		c.String(http.StatusBadRequest, "empty or invalid body")
		return
	}

	resp := make([]models.BatchShortURL, 0, len(req))
	userID := c.GetString("user_id")
	for _, item := range req {
		id, err := us.service.GenerateShortURL(item.OriginalURL, userID)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to create a short url")
			return
		}

		resp = append(resp, models.BatchShortURL{
			CorrelationID: item.CorrelationID,
			ShortURL:      us.baseURL + "/" + id,
		})
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(c.Writer).Encode(resp)
}

func (us *URLHandler) BatchDeleteHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}
	userID := userIDRaw.(string)

	var ids []string
	if err := c.ShouldBindJSON(&ids); err != nil || len(ids) == 0 {
		c.String(http.StatusBadRequest, "invalid request body")
		return
	}

	task := async.DeleteTask{
		UserID:    userID,
		ShortURLs: ids,
	}

	async.DeleteQueue <- task

	c.Writer.WriteHeader(http.StatusAccepted)
}

func (us *URLHandler) GetUserURLs(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}
	userID := userIDRaw.(string)

	urls, err := us.service.GetURLsByUser(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to get user urls")
		return
	}
	if len(urls) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	resp := make([]models.URLEntry, 0, len(urls))
	for _, u := range urls {
		resp = append(resp, models.URLEntry{
			ShortURL:    us.baseURL + "/" + u.ShortURL,
			OriginalURL: u.OriginalURL,
		})
	}

	c.JSON(http.StatusOK, resp)
}
