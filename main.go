package main

import (
	"api/cmd"
	"api/coind"
	"api/daemon"
	"api/db"
	_ "api/docs"
	"api/internal"
	"api/models"
	"api/services"
	"api/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

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
	app.Post("/submitTransaction", internal.SubmitTransaction)

	//swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	//external api
	app.Get("/api/inscriptions", getInscriptions)
	app.Get("/api/transactions", getTransaction)
	app.Post("/api/mint", mint)
	app.Post("/api/send", sendInscription)
	app.Get("/api/address", getAddress)

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

// Get new address godoc
// @Summary      Mint Inscription
// @Description  Mint Inscription
// @Tags         Daemon
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.NewAddressRequest
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /address [get]
func getAddress(c *fiber.Ctx) error {
	//dm := services.GetDaemon()
	type Address struct {
		Address string `json:"address"`
	}
	addr, err := cmd.CallJSONNonLock[Address]("bash", "-c", "/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300/wallet/ord --wallet ord wallet receive")
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	//privKey, err := coind.WrapDaemon(dm, 1, "dumpprivkey", addr.Address)
	//if err != nil {
	//	return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	//}
	return c.Status(http.StatusOK).JSON(&models.NewAddressRequest{
		Address: addr.Address,
	})
}

// Send inscription godoc
// @Summary      Send an Inscription
// @Description  Send an Inscription
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param 		 data body models.SendRequest true "File in base64 and file type"
// @Success      200  {object}  models.Inscribe
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /send [post]
func sendInscription(c *fiber.Ctx) error {
	var req models.SendRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}
	if req.Address == "" {
		return utils.ReportError(c, "Address is empty", http.StatusBadRequest)
	}
	if req.InscriptionID == "" {
		return utils.ReportError(c, "Inscription id is empty", http.StatusBadRequest)
	}
	dm := services.GetDaemon()
	sv, err := coind.WrapDaemon(dm, 1, "estimatesmartfee", 5, "economical")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}
	var fRate models.FeeRate
	err = json.Unmarshal(sv, &fRate)
	if err != nil {
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}

	feeRate := int(fRate.Feerate / 1024 * 100000000)

	if feeRate == 0 {
		return utils.ReportError(c, "FeeRate estimation error", http.StatusBadRequest)
	}
	if req.InscriptionID == "" {
		return utils.ReportError(c, "Inscription id is empty", http.StatusBadRequest)
	}

	s, err := cmd.CallJSON[models.Inscribe]("bash", "-c", fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300/wallet/ord --wallet ord wallet send --fee-rate %d %s %s", feeRate, req.Address, req.InscriptionID))
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(s)

}

// Mint godoc
// @Summary      Mint an Inscription
// @Description  Mint an Inscription
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param 		 data body models.MintRequest true "File in base64 and file type"
// @Success      200  {object}  models.Inscribe
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /mint [post]
func mint(c *fiber.Ctx) error {
	var req models.MintRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}

	if req.Format == "" {
		return utils.ReportError(c, "Format is empty", http.StatusBadRequest)
	}
	if req.Base64 == "" {
		return utils.ReportError(c, "Base64 is empty", http.StatusBadRequest)
	}

	dm := services.GetDaemon()
	sv, err := coind.WrapDaemon(dm, 1, "estimatesmartfee", 5, "economical")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}
	var fRate models.FeeRate
	err = json.Unmarshal(sv, &fRate)
	if err != nil {
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}

	feeRate := int(fRate.Feerate / 1024 * 100000000)

	byteArray, err := utils.DecodePayload([]byte(req.Base64))
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}
	fileType := strings.Split(req.Format, "/")[0]
	fileName := fmt.Sprintf("temp.%s", fileType)

	err = os.WriteFile(fileName, byteArray, 0644)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}

	s, err := cmd.CallJSON[models.Inscribe]("bash", "-c", fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300/wallet/ord --wallet ord inscribe --dry-run --fee-rate %d %s", feeRate, fileName))
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(s)

}

// Get detailed list of inscriptions godoc
// @Summary      List of Inscriptions in the wallet
// @Description  List of Inscriptions in the wallet
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success      200  {object}  models.ListInscriptionsResponse
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /inscriptions [get]
func getInscriptions(c *fiber.Ctx) error {
	pg := c.Query("page", "0")
	pgSize := c.Query("pageSize", "0")

	pgInt, err := strconv.Atoi(pg)
	pgSizeInt, err2 := strconv.Atoi(pgSize)
	if err != nil || err2 != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}
	req := &models.TxRequest{
		Page:     pgInt,
		PageSize: pgSizeInt,
	}
	pageSize := req.PageSize
	page := (req.Page - 1) * req.PageSize
	var res []models.TxTable
	var errDB error
	if pageSize == 0 && page == 0 {
		res, errDB = db.ReadArrayStruct[models.TxTable]("SELECT * FROM TRANSACTIONS_ORD")
	} else {
		res, errDB = db.ReadArrayStruct[models.TxTable](`SELECT * FROM TRANSACTIONS_ORD
WHERE oid NOT IN ( SELECT oid FROM TRANSACTIONS_ORD
                   ORDER BY id LIMIT ? )
ORDER BY id LIMIT ?`, page, pageSize)
	}
	if errDB != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	final := make([]models.TxTable, 0)
	for _, ins := range res {
		file := "./data_final/" + ins.OrdID[:8] + ".webp"
		bytes, err := utils.ReadFileAsBytes(file)
		if err != nil {
			return err
		}
		b64 := utils.EncodePayload(bytes)
		val := models.TxTable{
			ID:          ins.ID,
			OrdID:       ins.OrdID,
			TxID:        ins.TxID,
			Link:        ins.Link,
			ContentLink: ins.ContentLink,
			Base64:      b64,
		}
		final = append(final, val)
	}

	js := &models.ListInscriptionsResponse{
		HasError:     false,
		Status:       "OK",
		Inscriptions: final,
	}
	return c.Status(http.StatusOK).JSON(js)
}

// ListTransaction godoc
// @Summary      List of transactions in the BTC Core
// @Description  List of transactions in the BTC Core
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success      200  {object}  models.ListTransactions
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /transactions [get]
func getTransaction(c *fiber.Ctx) error {
	pg := c.Query("page", "0")
	pgSize := c.Query("pageSize", "0")

	pgInt, err := strconv.Atoi(pg)
	pgSizeInt, err2 := strconv.Atoi(pgSize)

	if err != nil || err2 != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}
	req := &models.TxRequest{
		Page:     pgInt,
		PageSize: pgSizeInt,
	}

	//checks
	if req.PageSize == 0 {
		pgSizeInt = 50
	}

	if req.Page == 0 {
		pgInt = 1
	}

	if req.PageSize < 1 && req.PageSize > 100 {
		return utils.ReportError(c, "Page size must be greater than 0 and not more than 100", http.StatusBadRequest)
	}

	pageSize := req.PageSize
	page := (req.Page - 1) * req.PageSize
	utils.ReportMessage(fmt.Sprintf("Offset: %d, Page: %d", page, req.PageSize))

	list, err := db.ReadArrayStruct[models.ListTransactionsDB](`SELECT * FROM LIST_TRANSACTIONS
WHERE oid NOT IN ( SELECT oid FROM LIST_TRANSACTIONS
                   ORDER BY id LIMIT ? )
ORDER BY id LIMIT ?`, page, pageSize)
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
