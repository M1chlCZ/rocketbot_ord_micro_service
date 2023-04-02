package main

import (
	"api/cmd"
	"api/coind"
	"api/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"api/db"
	"api/utils"
	"context"

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

	app := fiber.New(fiber.Config{
		AppName:       "Rocketbot ORD API",
		StrictRouting: false,
		WriteTimeout:  time.Second * 35,
		ReadTimeout:   time.Second * 35,
		IdleTimeout:   time.Second * 65,
	})

	utils.ReportMessage("Starting API, successfully authenticated")
	app.Post("/submitTransaction", submitTransaction)

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", 7500))
		if err != nil {
			utils.WrapErrorLog(err.Error())
			panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	utils.ReportMessage("/// = = Shutting down = = ///")
	defer cancel()
	_ = app.Shutdown()
	os.Exit(0)
}

func submitTransaction(c *fiber.Ctx) error {
	//txid := w.Get("txid")
	//coinid := w.Get("coinid")
	//cid, _ := strconv.Atoi(coinid)

	getInscriptions()
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		utils.STATUS: utils.OK,
	})
}

func getInscriptions() {
	s2, err := cmd.CallString("bash", "-c", "cat /home/dfwplay/.bitcoin/.cookie")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	s := strings.Split(s2, ":")
	daemon := &models.BitcoinDaemon{
		ID:         0,
		WalletUser: "__cookie__",
		WalletPass: s[1],
		WalletPort: 12300,
		Wallet:     "/wallet/ord",
		CoinID:     0,
		IP:         "127.0.0.1",
		PassPhrase: sql.NullString{},
	}

	callString, err := cmd.CallArrayJSON[models.Inscriptions]("bash", "-c", "/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300/wallet/ord --wallet ord wallet inscriptions")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	for _, ins := range callString {
		txid := strings.Split(ins.Location, ":")
		vout, err := strconv.Atoi(txid[1])
		if err != nil {
			utils.WrapErrorLog(err.Error())
			return
		}
		sv, err := coind.WrapDaemon(daemon, 1, "gettransaction", txid[0], false, true)
		if err != nil {
			utils.WrapErrorLog(err.Error())
			return
		}
		var info models.GetTransaction
		errJson := json.Unmarshal(sv, &info)
		if errJson != nil {
			utils.WrapErrorLog(errJson.Error())
			return
		}
		contentLink := fmt.Sprintf("https://ordinals.com/content/%s", ins.Inscription)
		addr := info.Decoded.Vout[vout].ScriptPubKey.Address
		_, err = db.InsertSQl(`INSERT INTO TRANSACTIONS_ORD (tx_id, ord_id, bc_address, link, content_link) 
									VALUES (?,?, ?, ?, ?)`, txid[0], ins.Inscription, addr, ins.Explorer, contentLink)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				utils.ReportMessage("Inscription already in db")
				continue
			}
			utils.WrapErrorLog(err.Error())
			continue
		}
	}

	in, err := db.ReadArrayStruct[models.TxTable]("SELECT * FROM TRANSACTIONS_ORD")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	marshal, err := json.Marshal(in)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}

	dst := &bytes.Buffer{}
	if err := json.Indent(dst, marshal, "", "  "); err != nil {
		panic(err)
	}

	utils.ReportSuccess(fmt.Sprintf("Inscriptions: %s", dst.String()))

	utils.ReportMessage("Inscriptions saved into db")

	time.Sleep(time.Second * 1)
}
