package controllers

import (
	"context"
	"fmt"
	"net/http"
	"project/initializers"
	models "project/model"
	services "project/service"
	"project/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	collection  *mongo.Collection
	ctx         context.Context
}

func NewAuthController(AuthService services.AuthService, UserService services.UserService, ctx context.Context, col *mongo.Collection) AuthController {
	return AuthController{
		authService: AuthService,
		userService: UserService,
		collection:  col,
		ctx:         ctx,
	}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var newUser *models.SignUpInput
	err := ctx.ShouldBindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if newUser.Password != newUser.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	code := utils.Encode(utils.GenerateRandomString(20))
	newUser.VerificationCode = code
	res, err := ac.authService.SignUpUser(newUser)
	if err != nil {
		ctx.JSON(http.StatusNotImplemented, gin.H{"message": "could not create the user", "error": err.Error()})
		return
	}

	//send email
	emailData := utils.EmailData{
		Code:    code,
		Subject: "Your accout verification code",
	}
	err = utils.SendEmail(res, &emailData)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "Cannot send Email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email with code sent"})
}

func (ac *AuthController) VerifyUser(ctx *gin.Context) {
	code := ctx.Params.ByName("verificationCode")
	err := ac.authService.VerifyUser(code)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": "Successfully Verified"})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Cannot Parse"})
		return
	}

	res, err := ac.userService.GetUserByMail(payload.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Problem with email"})
		return
	}

	fmt.Println(payload.Password, " ", res.Password)
	err = utils.VerifyPassword(payload.Password, res.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Incorrect Password"})
		return
	}

	if !res.Verified {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User Not Verified"})
		return
	}

	config, _ := initializers.LoadConfig(".")
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, res.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, res.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "could not refresh access token"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.GetUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	//Setting the MaxAge of a cookie to -1 means that the cookie will be deleted immediately when it reaches the client.
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Cannot Parse"})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	code := ctx.Params.ByName("resetCode")
	hashed_pw := utils.HashPassword(payload.Password)

	err := ac.authService.ResetPassword(hashed_pw, code)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})

}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload *models.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Cannot Parse"})
		return
	}

	res, err := ac.userService.GetUserByMail(payload.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Not Found"})
		return
	}

	if !res.Verified {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "UserNotVerified"})
		return
	}

	code := utils.Encode(utils.GenerateRandomString(20))
	err = ac.authService.ForgotPassword(res.Email, code)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	//send email
	emailData := utils.EmailData{
		Code:    code,
		Subject: "Your accout reser password code",
	}
	err = utils.SendEmail(res, &emailData)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "Cannot send Email"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "You will receive a reset email if user with that email exist"})
}
