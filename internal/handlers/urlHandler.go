package handlers

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rfruffer/go-musthave-shortener/internal/models"
	"github.com/rfruffer/go-musthave-shortener/internal/repository"
	"github.com/rfruffer/go-musthave-shortener/internal/services"
)

type URLHandler struct {
	service *services.URLService
	baseURL string
	db      repository.StoreRepositoryInterface
}

func NewURLHandler(service *services.URLService, baseURL string, db repository.StoreRepositoryInterface) *URLHandler {
	return &URLHandler{service: service, baseURL: baseURL, db: db}
}

func (us *URLHandler) CreateShortJSONURLHandler(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "empty or invalid body")
		return
	}

	id, err := us.service.GenerateShortURL(req.URL)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			existingShortID, getErr := us.db.GetShortIDByOriginalURL(req.URL)
			if getErr != nil {
				c.String(http.StatusInternalServerError, "internal error")
				return
			}
			resp := models.ShortenResponse{
				Result: us.baseURL + "/" + existingShortID,
			}
			c.JSON(http.StatusConflict, resp)
			return
		}
		c.String(http.StatusInternalServerError, "failed to create a short url")
		return
	}

	resp := models.ShortenResponse{
		Result: us.baseURL + "/" + id,
	}
	c.JSON(http.StatusCreated, resp)
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
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" { // 23505 = unique_violation
			existingShortID, getErr := us.db.GetShortIDByOriginalURL(originalURL)
			if getErr != nil {
				c.String(http.StatusInternalServerError, "internal error")
				return
			}
			c.Data(http.StatusConflict, "text/plain", []byte(us.baseURL+"/"+existingShortID))
			return
		}
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

func (us *URLHandler) Ping(c *gin.Context) {
	if err := us.db.Ping(); err != nil {
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

	for _, item := range req {
		id, err := us.service.GenerateShortURL(item.OriginalURL)
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
