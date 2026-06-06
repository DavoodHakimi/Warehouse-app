package router

import (
	"github.com/DavoodHakimi/warehouse-app/internal/auth"
	"github.com/DavoodHakimi/warehouse-app/internal/database"
	"github.com/DavoodHakimi/warehouse-app/internal/middleware"
	"github.com/DavoodHakimi/warehouse-app/internal/users"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	routerEng := gin.Default()
	v1 := routerEng.Group("/api/v1")
	db := database.GetDB()
	{
		authGroup := v1.Group("/auth")
		{
			authRepo := auth.NewRepository(db)
			authService := auth.NewService(authRepo)
			authHandler := auth.NewHandler(authService)

			authGroup.POST("/login", authHandler.LogInHandler)
			authGroup.POST("/signup", authHandler.SignUpHandler)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())

		{
			protected.GET("/auth/me", auth.MeHandler)

			usersGroup := protected.Group("/users")

			userRepo := users.NewRepository(db)
			userService := users.NewService(userRepo)
			userHandler := users.NewHandler(userService)
			{
				usersGroup.GET("/", userHandler.AllUsersHandler)             //returning list of company users, only ceo can see
				usersGroup.GET("/:userID", userHandler.UserInfoHandler)      //returning user details, only ceo can see
				usersGroup.POST("", userHandler.UserCreationHandler)         //creating new user, only ceo can do
				usersGroup.PATCH("/:userID", userHandler.UserUpdateHandler)  //only ceo and manager can edit user info, but not password
				usersGroup.DELETE("/:userID", userHandler.UserDeleteHandler) //only ceo can delete user
			}

			ordersGroup := protected.Group("/orders")
			{
				ordersGroup.GET("/")
				ordersGroup.GET("/:orderID")
				ordersGroup.POST("")
				ordersGroup.PATCH("/:orderID")
				ordersGroup.DELETE("/:orderID")
			}

			productsGroup := protected.Group("/products")
			{
				productsGroup.GET("/")
				productsGroup.POST("/")
				productsGroup.POST("/:productID/")
				productsGroup.PATCH("/:productID/")
				productsGroup.DELETE("/:productID/")
				productsGroup.POST("/:productID/approve")
				productsGroup.POST("/:productID/pack")
				productsGroup.POST("/:productID/ship")
				productsGroup.POST("/:productID/receive")

			}

			partnersGroup := protected.Group("/partners")
			{
				partnersGroup.GET("/")
				partnersGroup.GET("/:partnerID")
				partnersGroup.POST("")
				partnersGroup.PATCH("/:partnerID")
				partnersGroup.DELETE("/:partnerID")
			}
		}
	}

	return routerEng
}
