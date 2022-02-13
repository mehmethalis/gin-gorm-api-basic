package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"haliscicek.com/gin-api/dto"
	"haliscicek.com/gin-api/helper"
	"haliscicek.com/gin-api/model"
	"haliscicek.com/gin-api/service"
	"net/http"
	"strconv"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JwtService
}

func NewAuthController(authService service.AuthService, jwtService service.JwtService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDto
	err := ctx.ShouldBind(&loginDTO)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(model.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		tokenDetail, err := c.jwtService.ValidateToken(generatedToken)
		if err != nil {
			response := helper.BuildErrorResponse("Token can not validate", err.Error(), helper.EmptyObject{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		exp := tokenDetail.Claims.(jwt.MapClaims)["exp"].(float64)
		iat := tokenDetail.Claims.(jwt.MapClaims)["iat"].(float64)

		v.Token = generatedToken
		v.ExpiresAt = exp
		v.IssuedAt = iat
		response := helper.BuildResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Invalid credentials", "Please check credential", helper.EmptyObject{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDto
	err := ctx.ShouldBind(&registerDTO)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Dublicate email", helper.EmptyObject{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		tokenDetail, err := c.jwtService.ValidateToken(token)
		if err != nil {
			response := helper.BuildErrorResponse("Token can not validate", err.Error(), helper.EmptyObject{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		exp := tokenDetail.Claims.(jwt.MapClaims)["exp"].(float64)
		iat := tokenDetail.Claims.(jwt.MapClaims)["iat"].(float64)

		createdUser.ExpiresAt = exp
		createdUser.IssuedAt = iat
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}
