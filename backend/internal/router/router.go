package router

import (
	"github.com/DavoodHakimi/warehouse-app/internal/auth"
	"github.com/DavoodHakimi/warehouse-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	routerEng := gin.Default()
	v1 := routerEng.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", auth.LogInHandler)
			authGroup.POST("/signup", auth.SignUpHandler)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())

		{
			protected.GET("/auth/me", auth.MeHandler)

			usersGroup := protected.Group("/users")
			{
				usersGroup.GET("/")
				usersGroup.GET("/:userID")
				usersGroup.POST("")
				usersGroup.PATCH("/:userID")
				usersGroup.DELETE("/:userID")
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
