package server

import (
	"context"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/suyog1pathak/services/docs"
	"github.com/suyog1pathak/services/pkg/config"
	"github.com/suyog1pathak/services/pkg/controllers"
	"github.com/suyog1pathak/services/pkg/logger"
	log "github.com/suyog1pathak/services/pkg/logger"
	middlewarehealthcheck "github.com/suyog1pathak/services/pkg/middleware/healthcheck"
	middlewareservice "github.com/suyog1pathak/services/pkg/middleware/service"
	"github.com/suyog1pathak/services/pkg/model"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//	@title						services API
//	@version					0.1
//	@description				services API
//	@termsOfService				http://swagger.io/terms/
//	@BasePath					/
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func HandleRequest() {
	c := config.GetConfig()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	log.Info("starting server at", "port", c.App.ListeningPort)
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(c.App.ListeningPort),
		Handler: InitRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("error ", "error", err.Error())
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	stop()
	logger.Warn("shutting down gracefully, press Ctrl+C again to force", "draining_period", c.App.DrainingPeriod)

	// The context is used to inform the server it has DrainingPeriod seconds to finish
	// the request it is currently handling
	timeoutCtx, cancel := context.WithTimeout(context.Background(), c.App.DrainingPeriod*time.Second)
	defer cancel()

	go func() {
		if err := srv.Shutdown(timeoutCtx); err != nil {
			log.Error("error ", "error", err.Error())
		}
	}()

	// Wait for the timeout context to close.
	<-timeoutCtx.Done()

	if timeoutCtx.Err() == context.DeadlineExceeded {
		logger.Info("timeout exceeded, forcing shutdown")
	}

	logger.Info("Server exiting")
}

func InitRouter() *gin.Engine {
	docs.SwaggerInfo.Title = "services api"
	model.Setup()
	log := logger.Get()
	router := gin.New()
	router.Use(sloggin.New(log))
	router.Use(gin.Recovery())
	gin.SetMode(gin.ReleaseMode)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Info("please access swagger docs", "path", "http://localhost:8080/docs/index.html")
	{
		router.GET("/healthcheck", middlewarehealthcheck.HealthcheckCatchErrors(), controllers.Healthcheck)
		router.GET("/liveness", middlewarehealthcheck.HealthcheckCatchErrors(), controllers.Healthcheck)
		router.GET("/readiness", middlewarehealthcheck.HealthcheckCatchErrors(), controllers.Healthcheck)

		router.GET("/api/v1/services", middlewareservice.ServiceErrorHandler(), middlewareservice.ServiceQueryParams(), controllers.GetAllServices)
		router.GET("/api/v1/services/:name", middlewareservice.ServiceErrorHandler(), controllers.GetServiceByName)
		router.GET("/api/v1/services/:name/:version", middlewareservice.ServiceErrorHandler(), controllers.GetServiceNameAndVersion)
		router.POST("/api/v1/services", middlewareservice.ServiceErrorHandler(), middlewareservice.ServiceBodyValidation(), controllers.CreateService)
		router.PATCH("/api/v1/services/:name", middlewareservice.ServiceErrorHandler(), middlewareservice.ServiceBodyValidation(), controllers.UpdateService)
		router.PATCH("/api/v1/services/:name/:version", middlewareservice.ServiceErrorHandler(), middlewareservice.ServiceBodyValidation(), controllers.UpdateServiceVersion)
		router.DELETE("/api/v1/services/:name", middlewareservice.ServiceErrorHandler(), controllers.DeleteService)
		//v1.DELETE("/services/:name/:version", controllers.DeleteServiceVersion)
	}

	return router
}
