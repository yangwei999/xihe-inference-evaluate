package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/xihe-inference-evaluate/config"
	"github.com/opensourceways/xihe-inference-evaluate/controller"
	"github.com/opensourceways/xihe-inference-evaluate/docs"
	"github.com/opensourceways/xihe-inference-evaluate/domain/cloud"
	"github.com/opensourceways/xihe-inference-evaluate/domain/inference"
)

type Service struct {
	Inference inference.Inference
	Cloud     cloud.Cloud
}

func StartWebServer(port int, timeout time.Duration, cfg *config.Config, s *Service) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	setRouter(r, cfg, s)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

// setRouter init router
func setRouter(engine *gin.Engine, cfg *config.Config, s *Service) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "xihe"

	v1 := engine.Group(docs.SwaggerInfo.BasePath)
	{
		controller.AddRouterForInferenceController(
			v1, s.Inference,
		)

		controller.AddRouterForCloudController(
			v1, s.Cloud,
		)

	}

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		logrus.Infof(
			"| %d | %d | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
