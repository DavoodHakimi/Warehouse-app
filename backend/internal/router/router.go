package router

import (
	"github.com/DavoodHakimi/warehouse-app/internal/auth"
	"github.com/DavoodHakimi/warehouse-app/internal/middleware"
	"github.com/DavoodHakimi/warehouse-app/internal/orders"
	"github.com/DavoodHakimi/warehouse-app/internal/partners"
	"github.com/DavoodHakimi/warehouse-app/internal/products"
	"github.com/DavoodHakimi/warehouse-app/internal/users"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(db *gorm.DB) *gin.Engine {
	routerEng := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://frontend:3000"} // Your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true

	routerEng.Use(cors.New(config))

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
				orderRepo := orders.NewRepository(db)
				orderService := orders.NewService(orderRepo)
				orderHandler := orders.NewHandler(orderService)

				ordersGroup.GET("/", rbac.RBACMiddleware("orders.read"), orderHandler.AllOrdersHandler)
				ordersGroup.GET("/:orderID", rbac.RBACMiddleware("orders.read"), orderHandler.OrderInfoHandler)
				ordersGroup.POST("/", rbac.RBACMiddleware("orders.create"), orderHandler.OrderCreationHandler)
				ordersGroup.PATCH("/:orderID", rbac.RBACMiddleware("orders.update"), orderHandler.OrderUpdateHandler)
				ordersGroup.DELETE("/:orderID", rbac.RBACMiddleware("orders.delete"), orderHandler.OrderDeleteHandler)
				ordersGroup.POST("/:orderID/approve", rbac.RBACMiddleware("orders.update"), orderHandler.OrderApproveHandler)
				ordersGroup.POST("/:orderID/pack", rbac.RBACMiddleware("orders.pack"), orderHandler.OrderPackHandler)
				ordersGroup.POST("/:orderID/ship", rbac.RBACMiddleware("orders.ship"), orderHandler.OrderShipHandler)
				ordersGroup.POST("/:orderID/receive", rbac.RBACMiddleware("orders.receive"), orderHandler.OrderReceiveHandler)
				ordersGroup.POST("/:orderID/wait", rbac.RBACMiddleware("orders.update"), orderHandler.OrderWaitingHandler)
				ordersGroup.POST("/:orderID/cancel", rbac.RBACMiddleware("orders.update"), orderHandler.OrderCancelHandler)
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
