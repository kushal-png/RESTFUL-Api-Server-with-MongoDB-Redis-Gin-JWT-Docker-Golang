package controllers

import (
	"fmt"
	"net/http"
	models "project/model"
	services "project/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService services.PostServices
}

func NewPostController(p services.PostServices) PostController {
	return PostController{
		postService: p,
	}
}

func (pc *PostController) GetPost(ctx *gin.Context) {
	params := ctx.Params.ByName("postId")
	res, err := pc.postService.GetPost(params)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": res})
}

func (pc *PostController) GetPosts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	posts, err := pc.postService.GetPosts(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": posts})

}
func (pc *PostController) DeletePost(ctx *gin.Context) {
	params := ctx.Params.ByName("postId")
	err := pc.postService.DeletePost(params)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)
	payload := &models.CreatePost{
		User: currentUser.Email,
	}
	fmt.Println(payload.User)
	err := ctx.ShouldBindJSON(payload)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "errpr in parsing"})
		return
	}
    fmt.Println("enters")
	res, err := pc.postService.CreatePost(payload)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": res})

}
func (pc *PostController) UpdatePost(ctx *gin.Context) {
	params := ctx.Params.ByName("postId")
	var payload *models.UpdatePost
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "error in parsing"})
		return
	}
	res, err := pc.postService.UpdatePost(payload, params)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": res})

}
