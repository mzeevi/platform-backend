package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"

	"github.com/dana-team/platform-backend/src/auth"
	"github.com/dana-team/platform-backend/src/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine, tokenProvider auth.TokenProvider, scheme *runtime.Scheme) {
	v1 := engine.Group("/v1")

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	authGroup := v1.Group("/login")
	{
		authGroup.POST("", Login(tokenProvider))
	}

	namespacesGroup := v1.Group("/namespaces")
	namespacesGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	{
		namespacesGroup.GET("", GetNamespaces())
		namespacesGroup.GET("/:namespaceName", GetNamespace())
		namespacesGroup.POST("", CreateNamespace())
		namespacesGroup.DELETE("/:namespaceName", DeleteNamespace())
	}

	secretsGroup := namespacesGroup.Group("/:namespaceName/secrets")
	secretsGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	{
		secretsGroup.POST("", CreateSecret())
		secretsGroup.GET("", GetSecrets())
		secretsGroup.GET("/:secretName", GetSecret())
		secretsGroup.PUT("/:secretName", UpdateSecret())
		secretsGroup.DELETE("/:secretName", DeleteSecret())
	}

	cappGroup := namespacesGroup.Group("/:namespaceName/capps")
	cappGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	{
		cappGroup.POST("", CreateCapp())
		cappGroup.GET("", GetCapps())
		cappGroup.GET("/:cappName", GetCapp())
		cappGroup.PUT("/:cappName", UpdateCapp())
		cappGroup.PUT("/:cappName/state", EditCappState())
		cappGroup.DELETE("/:cappName", DeleteCapp())
	}

	cappRevisionGroup := namespacesGroup.Group("/:namespaceName/capprevisions")
	cappRevisionGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	{
		cappRevisionGroup.GET("", GetCappRevisions())
		cappRevisionGroup.GET("/:cappRevisionName", GetCappRevision())
	}

	usersGroup := namespacesGroup.Group("/:namespaceName/users")
	usersGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	{
		usersGroup.POST("", CreateUser())
		usersGroup.GET("", GetUsers())
		usersGroup.GET("/:userName", GetUser())
		usersGroup.PUT("/:userName", UpdateUser())
		usersGroup.DELETE("/:userName", DeleteUser())
	}

	configMapGroup := namespacesGroup.Group("/:namespaceName/configmaps")
	configMapGroup.Use(middleware.TokenAuthMiddleware(tokenProvider, scheme))
	{
		configMapGroup.GET("/:configMapName", GetConfigMap())
	}
}
