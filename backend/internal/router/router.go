package router

import (
	"github.com/DavoodHakimi/warehouse-app/internal/auth"
	"github.com/DavoodHakimi/warehouse-app/internal/middleware"
	"github.com/DavoodHakimi/warehouse-app/internal/partners"
	"github.com/DavoodHakimi/warehouse-app/internal/products"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func Setup(db *gorm.DB) *gin.Engine {
	routerEng := gin.Default()
	v1 := routerEng.Group("/api/v1")
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
		rbac := middleware.NewRBAC(db)

		{
			protected.GET("/auth/me", auth.MeHandler)

			usersGroup := protected.Group("/users")

			userRepo := users.NewRepository(db)
			userService := users.NewService(userRepo)
			userHandler := users.NewHandler(userService)
			{
				usersGroup.GET("/", rbac.RBACMiddleware("users.read"), userHandler.AllUsersHandler)
				usersGroup.POST("/", rbac.RBACMiddleware("users.create"), userHandler.UserCreationHandler)
				usersGroup.GET("/:userID", rbac.RBACMiddleware("users.read"), userHandler.UserInfoHandler)
				usersGroup.PATCH("/:userID", rbac.RBACMiddleware("users.update"), userHandler.UserUpdateHandler)
				usersGroup.DELETE("/:userID", rbac.RBACMiddleware("users.delete"), userHandler.UserDeleteHandler)
			}

			ordersGroup := protected.Group("/orders")
			{
				ordersGroup.GET("/", rbac.RBACMiddleware("orders.read"))
				ordersGroup.GET("/:orderID", rbac.RBACMiddleware("orders.read"))
				ordersGroup.POST("/", rbac.RBACMiddleware("orders.create"))
				ordersGroup.PATCH("/:orderID", rbac.RBACMiddleware("orders.update"))
				ordersGroup.DELETE("/:orderID", rbac.RBACMiddleware("orders.delete"))
				ordersGroup.POST("/:orderID/approve", rbac.RBACMiddleware("orders.update"))
				ordersGroup.POST("/:orderID/pack", rbac.RBACMiddleware("orders.pack"))
				ordersGroup.POST("/:orderID/ship", rbac.RBACMiddleware("orders.ship"))
				ordersGroup.POST("/:orderID/receive", rbac.RBACMiddleware("orders.receive"))
			}

			productsGroup := protected.Group("/products")
			{
				productRepo := products.NewRepository(db)
				productService := products.NewService(productRepo)
				productHandler := products.NewHandler(productService)

				productsGroup.GET("/", rbac.RBACMiddleware("products.read"), productHandler.AllProductsHandler)
				productsGroup.POST("/", rbac.RBACMiddleware("products.create"), productHandler.ProductCreationHandler)
				productsGroup.GET("/:productNumber/", rbac.RBACMiddleware("products.read"), productHandler.ProductInfoHandler)
				productsGroup.PATCH("/:productNumber/", rbac.RBACMiddleware("products.update"), productHandler.ProductUpdateHandler)
				productsGroup.DELETE("/:productNumber/", rbac.RBACMiddleware("products.delete"), productHandler.ProductDeleteHandler)

			}

			partnersGroup := protected.Group("/partners")
			{
				partnerRepo := partners.NewRepository(db)
				partnerService := partners.NewService(partnerRepo)
				partnerHandler := partners.NewHandler(partnerService)

				partnersGroup.GET("/", rbac.RBACMiddleware("partners.read"), partnerHandler.AllPartnersHandler)
				partnersGroup.POST("/", rbac.RBACMiddleware("partners.create"), partnerHandler.PartnerCreationHandler)
				partnersGroup.GET("/:PartnerID", rbac.RBACMiddleware("partners.read"), partnerHandler.PartnerInfoHandler)
				partnersGroup.PATCH("/:PartnerID", rbac.RBACMiddleware("partners.update"), partnerHandler.PartnerUpdateHandler)
				partnersGroup.DELETE("/:PartnerID", rbac.RBACMiddleware("partners.delete"), partnerHandler.PartnerDeleteHandler)
			}
		}
	}

	return routerEng
}
