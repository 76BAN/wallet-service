package main

import (
	"context"
	"fmt"
	"log"
	"wallets_service/database"
	"wallets_service/handler"
	middleware "wallets_service/midleware"
	"wallets_service/rabbitmq"
	"wallets_service/repository"
	routes "wallets_service/route"
	"wallets_service/service"

	"github.com/labstack/echo/v5"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Group Service berhasil terhubung ke database.")

	redisService := database.NewRedisService()

	err = redisService.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis berhasil terhubung")

	loggingPublisher, err := rabbitmq.NewLoggingPublisher(
		"amqp://alert_app:Silaban123@127.0.0.1:5672/",
	)
	if err != nil {
		log.Fatal("Gagal menghubungkan Logging Publisher: ", err)
	}
	defer loggingPublisher.Close()

	log.Println("RabbitMQ Logging Publisher berhasil terhubung")

	// Repository.
	walletRepository := repository.NewWalletsRepository(db)

	// Service.
	walletService := service.NewWalletService(walletRepository, loggingPublisher)

	// Handler.
	walletHandler := handler.NewWalletHandler(walletService)

	// JWT Middleware.
	jwtMiddleware, err := middleware.NewJWTMiddleware(
		"keys/public.pem",
		redisService,
	)
	if err != nil {
		log.Fatal("Gagal membuat JWT Middleware:", err)
	}

	//Membuat server Echo
	e := echo.New()

	// Endpoint untuk mengetes server
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Group Service berjalan",
		})
	})

	routes.Init(e, walletHandler, jwtMiddleware)

	//Menjalankan server pada port 8081
	log.Println("Group Service berjalan di http://localhost:8084")

	if err := e.Start(":8084"); err != nil {
		log.Fatal(err)
	}
}
