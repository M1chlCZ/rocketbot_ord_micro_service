package apis

import (
	"api/db"
	"api/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartLaunchpadApi() {
	db.InitMySQL()
	app := fiber.New(fiber.Config{
		AppName:       "Rocketbot Launchpad API",
		StrictRouting: false,
		WriteTimeout:  time.Second * 240,
		ReadTimeout:   time.Second * 240,
		IdleTimeout:   time.Second * 240,
	})
	app.Use(cors.New())

	app.Get("/launch/auth", getInscriptions)

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", 7700))
		if err != nil {
			utils.WrapErrorLog(err.Error())
			panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	//go getTransaction()
	utils.ReportMessage("<- Started Launchpad API ->")
	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	utils.ReportMessage("/// = = Shutting down = = ///")
	defer cancel()
	_ = app.Shutdown()
	os.Exit(0)

}
