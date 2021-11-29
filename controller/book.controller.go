package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"haliscicek.com/gin-api/dto"
	"haliscicek.com/gin-api/helper"
	"haliscicek.com/gin-api/model"
	"haliscicek.com/gin-api/service"
	"net/http"
	"strconv"
)

type BookController interface {
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JwtService
}

func NewBookController(bookService service.BookService, jwtService service.JwtService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) Insert(ctx *gin.Context) {
	var bookCreateDto dto.BookCreateDto
	authHeader := ctx.GetHeader("Authorization")
	userId := c.getUserIdByToken(authHeader)
	id, err := strconv.ParseUint(userId, 10, 64)
	if err == nil {
		bookCreateDto.UserID = id
	}
	errDto := ctx.ShouldBind(&bookCreateDto)
	if errDto != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDto.Error(), helper.EmptyObject{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result := c.bookService.Insert(bookCreateDto)
	res := helper.BuildResponse(true, "OK", result)
	ctx.JSON(http.StatusCreated, res)
}

func (c *bookController) Update(ctx *gin.Context) {
	var bookUpdateDto dto.BookUpdateDto

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := fmt.Sprintf("%v", claims["user_id"])
	bookId, err := strconv.ParseUint(ctx.Param("id"), 0, 0)

	if err != nil {
		res := helper.BuildErrorResponse("Failed to get id", "No param id were found", helper.EmptyObject{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if c.bookService.IsAllowedToEdit(userId, bookId) {
		id, errId := strconv.ParseUint(userId, 10, 64)
		if errId == nil {
			bookUpdateDto.UserID = id
			bookUpdateDto.ID = bookId
		}

		errDto := ctx.ShouldBind(&bookUpdateDto)
		if errDto != nil {
			res := helper.BuildErrorResponse("Failed to process request", errDto.Error(), helper.EmptyObject{})
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		result := c.bookService.Update(bookUpdateDto)
		res := helper.BuildResponse(true, "Ok", result)
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You dont have permission", "You are not owner", helper.EmptyObject{})
		ctx.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) Delete(ctx *gin.Context) {
	var book model.Book
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to get id", "No param id were found", helper.EmptyObject{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	book.ID = id
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userId, book.ID) {
		c.bookService.Delete(book)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObject{})
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You dont have permission", "You are not owner", helper.EmptyObject{})
		ctx.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) FindById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found ", err.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book model.Book = c.bookService.FindById(id)
	if (book == model.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObject{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	res := helper.BuildResponse(true, "Ok", book)
	ctx.JSON(http.StatusOK, res)

}

func (c *bookController) GetAll(ctx *gin.Context) {
	var books []model.Book = c.bookService.GetAll()
	res := helper.BuildResponse(true, "ok!", books)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) getUserIdByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}
