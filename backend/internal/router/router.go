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
				usersGroup.GET("/", middleware.RBACMiddleware("users.read"), userHandler.AllUsersHandler)
				usersGroup.POST("/", middleware.RBACMiddleware("users.create"), userHandler.UserCreationHandler)
				usersGroup.GET("/:userID", middleware.RBACMiddleware("users.read"), userHandler.UserInfoHandler)
				usersGroup.PATCH("/:userID", middleware.RBACMiddleware("users.update"), userHandler.UserUpdateHandler)
				usersGroup.DELETE("/:userID", middleware.RBACMiddleware("users.delete"), userHandler.UserDeleteHandler)
			}

			ordersGroup := protected.Group("/orders")
			{
				ordersGroup.GET("/", middleware.RBACMiddleware("orders.read"))
				ordersGroup.GET("/:orderID", middleware.RBACMiddleware("orders.read"))
				ordersGroup.POST("/", middleware.RBACMiddleware("orders.create"))
				ordersGroup.PATCH("/:orderID", middleware.RBACMiddleware("orders.update"))
				ordersGroup.DELETE("/:orderID", middleware.RBACMiddleware("orders.delete"))
				ordersGroup.POST("/:orderID/approve", middleware.RBACMiddleware("orders.update"))
				ordersGroup.POST("/:orderID/pack", middleware.RBACMiddleware("orders.pack"))
				ordersGroup.POST("/:orderID/ship", middleware.RBACMiddleware("orders.ship"))
				ordersGroup.POST("/:orderID/receive", middleware.RBACMiddleware("orders.receive"))
			}

			productsGroup := protected.Group("/products")
			{
				productRepo := products.NewRepository(db)
				productService := products.NewService(productRepo)
				productHandler := products.NewHandler(productService)

				productsGroup.GET("/", middleware.RBACMiddleware("products.read"), productHandler.AllProductsHandler)
				productsGroup.POST("/", middleware.RBACMiddleware("products.create"), productHandler.ProductCreationHandler)
				productsGroup.GET("/:productNumber/", middleware.RBACMiddleware("products.read"), productHandler.ProductInfoHandler)
				productsGroup.PATCH("/:productNumber/", middleware.RBACMiddleware("products.update"), productHandler.ProductUpdateHandler)
				productsGroup.DELETE("/:productNumber/", middleware.RBACMiddleware("products.delete"), productHandler.ProductDeleteHandler)

			}

			partnersGroup := protected.Group("/partners")
			{
				partnerRepo := partners.NewRepository(db)
				partnerService := partners.NewService(partnerRepo)
				partnerHandler := partners.NewHandler(partnerService)

				partnersGroup.GET("/", middleware.RBACMiddleware("partners.read"), partnerHandler.AllPartnersHandler)
				partnersGroup.POST("/", middleware.RBACMiddleware("partners.create"), partnerHandler.PartnerCreationHandler)
				partnersGroup.GET("/:PartnerID", middleware.RBACMiddleware("partners.read"), partnerHandler.PartnerInfoHandler)
				partnersGroup.PATCH("/:PartnerID", middleware.RBACMiddleware("partners.update"), partnerHandler.PartnerUpdateHandler)
				partnersGroup.DELETE("/:PartnerID", middleware.RBACMiddleware("partners.delete"), partnerHandler.PartnerDeleteHandler)
			}
		}
	}

	return routerEng
}
