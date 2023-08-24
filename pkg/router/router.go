package router

import (
	"market/pkg/handlers/segment"
	"market/pkg/handlers/user"
	"market/pkg/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repositories repository.IRepositories) *gin.Engine {
	router := gin.Default()

	usersH := user.NewHandler(repositories.GetUsersRepo())
	segmentsH := segment.NewHandler(repositories.GetSegmentsRepo())

	v1 := router.Group("/api/v1")
	{
		usersG := v1.Group("/users")
		{
			usersG.GET("", usersH.GetUsersSlugs)
			usersG.PATCH("", usersH.EditUsersSlugs)
		}

		segmentsG := v1.Group("/segments")
		{
			segmentsG.POST("", segmentsH.CreateSegments)
			segmentsG.DELETE("", segmentsH.DeleteSegments)
		}

		v1.StaticFile("/swagger/api.json", "./api/api.json")
		v1.Static("/swagger-ui", "./static/swagger-ui/dist")
	}

	return router
}
