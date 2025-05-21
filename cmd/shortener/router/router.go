package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rfruffer/go-musthave-shortener/internal/handlers"
	"github.com/rfruffer/go-musthave-shortener/internal/middlewares"
	"go.uber.org/zap"
)

type Router struct {
	URLHandler *handlers.URLHandler
}

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

	r.POST("/", rt.URLHandler.CreateShortURLHandler)
	r.GET("/:id", rt.URLHandler.GetShortURLHandler)

	api := r.Group("/api")
	api.Use(middlewares.GinGzipMiddleware())
	api.POST("/shorten", rt.URLHandler.CreateShortJSONURLHandler)

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "invalid request")
	})

	r.NoMethod(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "invalid request")
	})

	return r
}
