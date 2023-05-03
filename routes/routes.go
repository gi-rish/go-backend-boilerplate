package routes

import (
	"fmt"
	"time"

	"github.com/NewStreetTechnologies/go-backend-boilerplate/controllers"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/middleware"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouterConfig struct {
	R                     *gin.Engine
	TranslationRepository repository.TranslationRepository
	Logger                *logrus.Logger
}

func (c *RouterConfig) InitRouter() *gin.Engine {
	if gin.Mode() != gin.TestMode {
		c.R.Use(gin.Recovery())
		c.R.Use(middleware.SampleMiddleware)
		c.R.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] \"%s - %s %s %s %d %s \"%s\" %s\"\n",
				param.TimeStamp.Format(time.RFC1123),
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		}))
	}

	c.customRoutes()

	return c.R
}

func (c *RouterConfig) customRoutes() {
	v1 := c.R.Group("/api/v1")
	{
		translationRouter := v1.Group("/translations")
		{
			translations := &controllers.TranslationController{
				TranslationRepository: c.TranslationRepository,
				Logger:                c.Logger,
			}
			translationRouter.POST("/create", translations.CreateTranslation)
			translationRouter.GET("/:id", translations.GetTranslationsByLanguage)
			translationRouter.GET("/", translations.GetDefaultLanguageTranslations)
			translationRouter.GET("/all", translations.GetAllTranslations)
		}

	}
}
