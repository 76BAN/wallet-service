package routes

import (
	"wallets_service/handler"
	middleware "wallets_service/midleware"

	"github.com/labstack/echo/v5"
)

func Init(
	e *echo.Echo,
	walletHandler *handler.WalletHandler,
	midleware *middleware.JWTMiddleware,

) {

	e.POST("/wallets", walletHandler.CreateWallet, midleware.Authenticate)
	e.GET("/wallets/:id", walletHandler.FindWalletByID, midleware.Authenticate)
	e.GET("/group/wallets/:groupId", walletHandler.FindWalletsByGroupID, midleware.Authenticate)
	e.GET("/user/group/wallets/:groupId", walletHandler.FindWalletsByUserIDAndGroupID, midleware.Authenticate)
	e.DELETE("/wallets/:id", walletHandler.DeleteWallet, midleware.Authenticate)
	e.PUT("/wallets/:id", walletHandler.UpdateWallet, midleware.Authenticate)
}
