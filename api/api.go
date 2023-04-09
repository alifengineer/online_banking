package api

import (
	"github.com/dilmurodov/online_banking/api/docs"
	"github.com/dilmurodov/online_banking/api/handlers"
	"github.com/dilmurodov/online_banking/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetUpRouter godoc
// @description This is online banking API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetUpRouter(h handlers.Handler, cfg config.Config) (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	docs.SwaggerInfo.Title = cfg.ServiceName
	docs.SwaggerInfo.Version = cfg.Version
	docs.SwaggerInfo.Schemes = []string{cfg.HTTPScheme}

	r.Use(customCORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// auth
			auth := v1.Group("/auth")
			{
				// авторизация
				auth.POST("/login", h.LoginHandler)
				// регистрация
				auth.POST("/register", h.RegisterHandler)
			}

			// user
			account := v1.Group("/user")
			account.Use(h.AuthMiddleware)
			{
				// создание счета
				account.POST("/accounts", h.AccountCreateHandler)
				// получение списка счетов
				account.GET("/accounts", h.AccountsGetHandler)
				// получение счета по id
				account.GET("/accounts/:id", h.AccountGetHandler)
				// получение транзакций по счету
				account.GET("/accounts/:id/transactions", h.AccountTransactionsHandler)
				// получение транзакции по id
				account.GET("/accounts/:id/transactions/:transaction_id", h.AccountTransactionByIDHandler)
			}

			// payments
			payments := v1.Group("/payments")
			payments.Use(h.AuthMiddleware)
			{
				// вывод денег
				payments.POST("/withdrawal", h.WithDrawalHandler)
				// пополнение счета
				payments.POST("/deposit", h.DepositHandler)
				// перевод денег
				payments.POST("authorization", h.ConfirmPaymentHandler)
				// подтверждение перевода
				payments.POST("/capture", h.CaptureTransactionsHandler)
				// перевод на чужой счет
				payments.POST("/transfer", h.TransferHandler)
			}
		}
	}
	return
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
