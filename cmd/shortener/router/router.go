package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/middlewares"
	"go.uber.org/zap"
)

// Router предоставляет методы для работы с handlers.
type Router struct {
	URLHandler *handlers.URLHandler
	SecretKey  string
}

// SetupRouter устанавливает маршруты для handlers.
func SetupRouter(rt Router) http.Handler {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	r := gin.New()

	middlewares.InitLogger(sugar)
	r.Use(middlewares.GinLoggingMiddleware())
	r.Use(gin.Recovery())
	r.Use(middlewares.AuthMiddleware(rt.SecretKey))

	r.POST("/", rt.URLHandler.CreateShortURLHandler)
	r.GET("/:id", rt.URLHandler.GetShortURLHandler)
	r.GET("/ping", rt.URLHandler.Ping)

	api := r.Group("/api")
	api.Use(middlewares.GinGzipMiddleware())
	api.POST("/shorten", rt.URLHandler.CreateShortJSONURLHandler)
	api.POST("/shorten/batch", rt.URLHandler.Batch)
	api.GET("/user/urls", rt.URLHandler.GetUserURLs)
	api.DELETE("/user/urls", rt.URLHandler.BatchDeleteHandler)

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "invalid request")
	})

	r.NoMethod(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "invalid request")
	})

	return r
}
