package main

import (
	"github.com/gin-contrib/cors"
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
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	jwtService     service.JwtService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	bookServie     service.BookService       = service.NewBookService(bookRepository)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookServie, jwtService)
)

func main() {
	defer config.CloseDbConnection(db)
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Authorization"," Content-Type"}
	r.Use(cors.New(corsConfig))

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

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.GetAll)
		bookRoutes.GET("/:id", bookController.FindById)
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}
	r.Run()
}
