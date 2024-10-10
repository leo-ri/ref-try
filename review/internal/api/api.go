package api

import (
	// "context"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host string
	Port int
}

type API struct {
	srv *http.Server
	MUX *gin.Engine
	svc Service
	CFG Config
}

func New(cfg Config, svc Service) API {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	return API{
		MUX: r,
		CFG: cfg,
		svc: svc,
	}.withServer().withRoutes()
}

func (a API) withServer() API {
	a.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.CFG.Port),
		Handler: a.MUX,
	}
	return a
}

func (a API) withRoutes() API { // TODO ??
	apiGroup := a.MUX.Group("/api")
	apiGroup.POST("/apply", a.Apply)
	apiGroup.POST("/create", a.Create)
	apiGroup.GET("/coupons", a.Get)
	return a
}

func (a API) Start() error {
	if err := a.srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (a API) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server shutted down")
	return nil
}
