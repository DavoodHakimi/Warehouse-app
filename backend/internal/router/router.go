package router

import (
	"github.com/DavoodHakimi/warehouse-app/internal/auth"
	"github.com/DavoodHakimi/warehouse-app/internal/database"
	"github.com/DavoodHakimi/warehouse-app/internal/middleware"
	"github.com/DavoodHakimi/warehouse-app/internal/partners"
	"github.com/DavoodHakimi/warehouse-app/internal/products"
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
				usersGroup.GET("/", userHandler.AllUsersHandler)
				usersGroup.GET("/:userID", userHandler.UserInfoHandler)
				usersGroup.POST("", userHandler.UserCreationHandler)
				usersGroup.PATCH("/:userID", userHandler.UserUpdateHandler)
				usersGroup.DELETE("/:userID", userHandler.UserDeleteHandler)
			}

			ordersGroup := protected.Group("/orders")
			{
				ordersGroup.GET("/")
				ordersGroup.GET("/:orderID")
				ordersGroup.POST("")
				ordersGroup.PATCH("/:orderID")
				ordersGroup.DELETE("/:orderID")
				ordersGroup.POST("/:orderID/approve")
				ordersGroup.POST("/:orderID/pack")
				ordersGroup.POST("/:orderID/ship")
				ordersGroup.POST("/:orderID/receive")
			}

			productsGroup := protected.Group("/products")
			{
				productRepo := products.NewRepository(db)
				productService := products.NewService(productRepo)
				productHandler := products.NewHandler(productService)

				productsGroup.GET("/", productHandler.AllProductsHandler)
				productsGroup.POST("/", productHandler.ProductCreationHandler)
				productsGroup.GET("/:productNumber/", productHandler.ProductInfoHandler)
				productsGroup.PATCH("/:productNumber/", productHandler.ProductUpdateHandler)
				productsGroup.DELETE("/:productNumber/", productHandler.ProductDeleteHandler)

			}

			partnersGroup := protected.Group("/partners")
			{
				partnerRepo := partners.NewRepository(db)
				partnerService := partners.NewService(partnerRepo)
				partnerHandler := partners.NewHandler(partnerService)

				partnersGroup.GET("/", partnerHandler.AllPartnersHandler)
				partnersGroup.GET("/:PartnerID", partnerHandler.PartnerInfoHandler)
				partnersGroup.POST("", partnerHandler.PartnerCreationHandler)
				partnersGroup.PATCH("/:PartnerID", partnerHandler.PartnerUpdateHandler)
				partnersGroup.DELETE("/:PartnerID", partnerHandler.PartnerDeleteHandler)
			}
		}
	}

	return routerEng
}
