package routes

import (
	"project/controllers"
	"project/middleware"
	services "project/service"

	"github.com/gin-gonic/gin"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewPostRouteController(controller controllers.PostController) PostRouteController {
	return PostRouteController{postController: controller}
}

func (pc *PostRouteController) PostRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/posts")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("", pc.postController.GetPosts)
	router.GET("/:postId", pc.postController.GetPost)
	router.DELETE("/:postId", pc.postController.DeletePost)
	router.POST("/post", pc.postController.CreatePost)
	router.PATCH("/:postId", pc.postController.UpdatePost)
}
