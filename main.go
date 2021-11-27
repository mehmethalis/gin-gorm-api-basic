package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"haliscicek.com/gin-api/config"
	"haliscicek.com/gin-api/controller"
	"haliscicek.com/gin-api/middleware"
	"haliscicek.com/gin-api/repository"
	"haliscicek.com/gin-api/service"
)

var (
	db             *gorm.DB                  = config.DbConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	jwtService     service.JwtService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
)

func main() {
	defer config.CloseDbConnection(db)
	r := gin.Default()
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)

	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}
	r.Run()
}
