package main

import (
	"api/daemon"
	"api/db"
	_ "api/docs"
	"api/models"
	"api/services"
	"api/utils"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"
)

// @title Rocketbot ORD API
// @version 1.0
// @description Private API for ORD
// @termsOfService http://swagger.io/terms/

// @contact.name RocketBot
// @contact.url http://app.rocketbot.pro
// @contact.email m1chlcz18@gmail.com
// @contact.name Michal Žídek

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host     89.116.25.234:7500
// @BasePath /api
func main() {
	err := db.InitDB()
	if err != nil {
		utils.WrapErrorLog(fmt.Sprintf("Error opening db: %s", err.Error()))
		return
	}
	go services.GetInscriptions()
	go daemon.StartCron()

	app := fiber.New(fiber.Config{
		AppName:       "Rocketbot ORD API",
		StrictRouting: false,
		WriteTimeout:  time.Second * 240,
		ReadTimeout:   time.Second * 240,
		IdleTimeout:   time.Second * 240,
	})
	app.Use(cors.New())
	//internal api
	app.Post("/submitTransaction", submitTransaction)

	//swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	//external api
	app.Get("/api/getInscriptions", getInscriptions)
	app.Post("/api/getTransactions", getTransaction)

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", 7500))
		if err != nil {
			utils.WrapErrorLog(err.Error())
			panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	//go getTransaction()
	utils.ReportMessage("<- Started ORD API ->")
	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	utils.ReportMessage("/// = = Shutting down = = ///")
	defer cancel()
	_ = app.Shutdown()
	os.Exit(0)
}

func submitTransaction(c *fiber.Ctx) error {
	//curl -X POST -H "txid:$1" -H "coinid:$coinID" -H "node_id:$nodeID" http://localhost:7466/submitTransaction
	txid := c.Get("txid")

	mp := &fiber.Map{
		"txid":   txid,
		"coinid": 0,
	}

	go services.GetInscriptions()
	go services.SaveListTransaction()

	for {
		r, err := utils.POSTRequest[models.ErrorHTTP]("submitTransaction", mp)
		if err != nil {
			utils.WrapErrorLog(err.Error())
			time.Sleep(time.Second * 5)
			continue
		}
		if r.HasError != false {
			utils.WrapErrorLog(r.ErrorMessage)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		utils.STATUS: utils.OK,
	})
}

func getInscriptions(c *fiber.Ctx) error {
	res, err := db.ReadArray[models.TxTable]("SELECT * FROM TRANSACTIONS_ORD")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			utils.STATUS: utils.OK,
		})
	}
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		utils.STATUS: utils.OK,
		utils.ERROR:  false,
		"txs":        res,
	})
}

// ListTransaction godoc
// @Summary      List transactions from BTC Core
// @Description  List transactions from BTC Core
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param 		 data body models.TxRequest true "Info about device"
// @Success      200  {object}  models.ListTransactions
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /getTransactions [post]
func getTransaction(c *fiber.Ctx) error {

	var req models.TxRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}

	pageSize := req.PageSize
	offset := (req.Page - 1) * pageSize

	list, err := db.ReadArrayStruct[models.ListTransactionsDB](`SELECT * FROM LIST_TRANSACTIONS LIMIT ?, ?`, pageSize, offset)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	tx := make([]models.ListTransactionsDB, 0)
	for _, v := range list {
		tx = append(tx, v)
	}
	//marshal, err := json.Marshal(tx)
	//if err != nil {
	//	utils.WrapErrorLog(err.Error())
	//	return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	//}
	return c.JSON(tx)
}
