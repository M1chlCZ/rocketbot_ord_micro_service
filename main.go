package main

import (
	"api/cmd"
	"api/coind"
	"api/daemon"
	"api/db"
	"api/models"
	"api/services"
	"api/utils"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"

	"syscall"
	"time"
)

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
	//internal api
	app.Post("/submitTransaction", submitTransaction)

	//external api
	app.Get("/getInscriptions", getInscriptions)

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
	//TODO get transactions from bitcoin daemon
	r, err := utils.POSTRequest[models.HttpErr]("http://localhost:7466/submitTransaction", mp)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			utils.STATUS: utils.OK,
		})
	}
	if r.HasError != false {
		utils.WrapErrorLog(r.Message)
		return c.Status(http.StatusOK).JSON(&fiber.Map{
			utils.STATUS: utils.OK,
		})
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

func saveListTransaction() {
	s2, err := cmd.CallString("bash", "-c", "cat /home/dfwplay/.bitcoin/.cookie")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}

	s := strings.Split(s2, ":")
	dm := &models.BitcoinDaemon{
		ID:         0,
		WalletUser: "__cookie__",
		WalletPass: s[1],
		WalletPort: 12300,
		Wallet:     "/wallet/ord",
		CoinID:     0,
		IP:         "127.0.0.1",
		PassPhrase: sql.NullString{},
	}
	sv, err := coind.WrapDaemon(dm, 1, "listtransactions", "*", 999999)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	var model models.ListTransactions
	err = json.Unmarshal(sv, &model)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	for _, v := range model {
		_, err := db.InsertSQl(`INSERT INTO LIST_TRANSACTIONS (address, category, amount, vout, fee, confirmations, blockhash, blockheight, blockindex, blocktime, txid, wtxid, time, timereceived, 
                               bip125_replaceable, abandoned, label, trusted) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			v.Address, v.Category, v.Amount, v.Vout, v.Fee, v.Confirmations, v.Blockhash, v.Blockheight, v.Blockindex, v.Blocktime, v.Txid, v.Wtxid, v.Time, v.Timereceived, v.Bip125Replaceable, v.Abandoned, v.Label, v.Trusted)
		if err != nil {
			if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
				utils.WrapErrorLog(err.Error())
			}
			continue
		}
	}
}

func getTransaction() {
	list, err := db.ReadArrayStruct[models.ListTransactionsDB]("SELECT * FROM LIST_TRANSACTIONS")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	tx := make([]models.ListTransactionsDB, 0)
	for _, v := range list {
		tx = append(tx, v)
	}
	marshal, err := json.Marshal(tx)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, marshal, "", "  "); err != nil {
		panic(err)
	}
	utils.ReportMessage(dst.String())
}
