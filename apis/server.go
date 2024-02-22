package apis

import (
	"api/cmd"
	"api/coind"
	"api/daemon"
	"api/db"
	_ "api/docs"
	"api/grpcClient"
	"api/grpcModels"
	"api/internal"
	"api/models"
	"api/services"
	"api/utils"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func StartORDApi() {
	//err := gen.New().Build(&gen.Config{
	//	SearchDir:          "./",
	//	MainAPIFile:        "./apis/server.go",
	//	PropNamingStrategy: swag.CamelCase,
	//	OutputDir:          "./docs",
	//	OutputTypes:        []string{"json", "yaml"},
	//	ParseDependency:    true,
	//	// This is a swag/v2 field
	//	GenerateOpenAPI3Doc: false,
	//	ParseVendor:         true,
	//})
	//if err != nil {
	//	log.Fatalf("Failed to generate swagger.json err: %v", err)
	//}
	err := db.InitDB()
	if err != nil {
		utils.WrapErrorLog(fmt.Sprintf("Error opening db: %s", err.Error()))
		return
	}
	go services.GetInscriptions()
	go daemon.StartCron()
	go func() {
		_, err := grpcClient.PingMainNode(&grpcModels.PingRequest{Ping: true})
		if err != nil {
			utils.WrapErrorLog(err.Error())
			return
		}
	}()

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
	app.Get("/api/inscriptions/nsfw", getNsfwInscription)
	app.Get("/api/inscriptions/nsfw/approve", approveNsfwInscription)
	app.Get("/api/inscription/image", getImage)
	app.Get("/api/transactions", getTransaction)
	app.Get("/api/transaction/raw", getRawTx)
	app.Post("/api/mint", mint)
	app.Post("/api/inscription", createInscription)
	app.Post("/api/estimate", estimate)
	app.Post("/api/send", sendInscription)
	app.Get("/api/address", getAddress)
	app.Get("/api/feerate", getFeeRate)
	app.Post("/api/fee/quality/estimate", EstimateFeeQuality)

	app.Post("api/nsfw", testPic)

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", 7500))
		if err != nil {
			utils.WrapErrorLog(err.Error())
			panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	utils.ReportMessage("<- Started ORD API ->")
	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	utils.ReportMessage("/// = = Shutting down = = ///")
	defer cancel()
	_ = app.Shutdown()
	os.Exit(0)
}

func createInscription(c *fiber.Ctx) error {
	var req models.MintLaunchRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}
	utils.ReportMessage(fmt.Sprintf("Mint inscription id"))
	if req.Format == "" {
		return utils.ReportError(c, "Format is empty", http.StatusBadRequest)
	}
	if req.Base64 == "" {
		return utils.ReportError(c, "Base64 is empty", http.StatusBadRequest)
	}

	if req.FeeRate == 0 {
		return utils.ReportError(c, "Fee rate is empty", http.StatusBadRequest)
	}

	if req.Addr == "" {
		return utils.ReportError(c, "Address is empty", http.StatusBadRequest)
	}

	//dm := services.GetDaemon()
	//sv, err := coind.WrapDaemon(dm, 1, "estimatesmartfee", 5, "economical")
	//if err != nil {
	//	utils.WrapErrorLog(err.Error())
	//	return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	//}
	//var fRate models.FeeRate
	//err = json.Unmarshal(sv, &fRate)
	//if err != nil {
	//	return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	//}
	//
	//feeRate := int(fRate.Feerate / 1024 * 100000000)
	//fileType := strings.Split(req.Format, "/")[0]
	filename := fmt.Sprintf("/home/dfwplay/api/data_ins/%s.%s", time.Now().Format("200601021504"), req.Format)
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return utils.ReportError(c, "Can't create folder", http.StatusBadRequest)
		}
	}
	// Check if file exists
	if _, err := os.Stat(filename); err == nil {
		return utils.ReportError(c, "Something stat", http.StatusBadRequest)
	}

	file, err := os.Create(filename)
	if err != nil {
		return utils.ReportError(c, "Can't create file", http.StatusBadRequest)
	}

	fl, err := utils.DecodePayload([]byte(req.Base64))
	if err != nil {
		return utils.ReportError(c, "Can't decode payload", http.StatusBadRequest)
	}

	_, err = io.Copy(file, bytes.NewReader(fl))
	if err != nil {
		return utils.ReportError(c, "Can't copy bytes", http.StatusBadRequest)
	}

	//tx := &grpcModels.NSFWRequest{
	//    Base64:   req.Base64,
	//    Filename: fmt.Sprintf("pic.%s", req.Format),
	//}
	//res, err := grpcClient.DetectNSFW(tx)
	//if err != nil {
	//    utils.WrapErrorLog(err.Error())
	//    return utils.ReportError(c, "Cannot check if NSFW image", http.StatusInternalServerError)
	//}
	//
	//if res.NsfwPicture {
	//    err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", filename))).Run()
	//    if err != nil {
	//        utils.WrapErrorLog("Can't delete file in data")
	//    }
	//    return utils.ReportError(c, "NSFW image", http.StatusConflict)
	//}
	//
	//if res.NsfwText {
	//    err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", fileName))).Run()
	//    if err != nil {
	//        utils.WrapErrorLog("Can't delete file in data")
	//    }
	//    return utils.ReportError(c, "NSFW Text in the image", http.StatusConflict)
	//}
	go inscribe(req.ID, req.FeeRate, filename, req.Addr)

	return c.Status(200).JSON(&fiber.Map{
		utils.STATUS: utils.OK,
		utils.ERROR:  false,
	})
}

var inscriptionMutex sync.Mutex

func inscribe(id int64, feeRate int, filename, destination string) {
	inscriptionMutex.Lock()
	defer inscriptionMutex.Unlock()
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			utils.WrapErrorLog(err.Error())
		}
	}(filename)
	utils.ReportMessage(fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file /home/dfwplay/.bitcoin/.cookie --rpc-url 127.0.0.1:12300 wallet inscribe --destination %s --fee-rate %d %s", destination, feeRate, filename))
	s, err := cmd.CallJSON[models.Inscribe]("bash", "-c", fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file /home/dfwplay/.bitcoin/.cookie --rpc-url 127.0.0.1:12300  wallet inscribe --destination %s --fee-rate %d %s", destination, feeRate, filename))
	if err != nil {
		utils.WrapErrorLog(err.Error())
		_, errWeb := utils.POSTRequest[models.ErrorHTTP]("https://rocket.art/api/api/v1/inscription/error", &fiber.Map{
			"error": err.Error(),
			"id":    id,
		})
		if errWeb != nil {
			utils.WrapErrorLog(errWeb.Error())
		}
		return

	}
	utils.ReportMessage(fmt.Sprintf("Inscription created: %v", s))
	_, errWeb := utils.POSTRequest[models.ErrorHTTP]("https://rocket.art/api/api/v1/inscription/confirm", &fiber.Map{
		"commit":      s.Commit,
		"inscription": s.Inscription,
		"reveal":      s.Reveal,
		"fees":        s.Fees,
		"id":          id,
	})
	if errWeb != nil {
		utils.WrapErrorLog(errWeb.Error())
	}
	return
}

// Approve NSFW inscription
// @Summary      Approve Inscription from NSFW list
// @Description  Approve Inscription from NSFW list
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param ord query string true "ORD id"
// @Success      200  {object}  models.HttpSuccess
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /inscriptions/nsfw/approve [get]
func approveNsfwInscription(c *fiber.Ctx) error {
	pg := c.Query("ord", "")
	if pg == "" {
		return utils.ReportError(c, "No inscription provided", http.StatusBadRequest)
	}
	_, err := db.InsertSQl("UPDATE NSFW_ORD SET approved = 1 WHERE ord_id = ?", pg)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	nsfw, err := db.ReadStruct[models.NSFWTable]("SELECT * FROM NSFW_ORD WHERE ord_id = ?", pg)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	_, err = db.InsertSQl(`INSERT INTO TRANSACTIONS_ORD (tx_id, ord_id, bc_address, link, content_link) 
									VALUES (?,?, ?, ?, ?)`, nsfw.TxID, nsfw.OrdID, nsfw.BcAddress, nsfw.Link, nsfw.ContentLink)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		//return utils.ReportError(c, "Can't create approved ordinal", http.StatusInternalServerError)
	}
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		utils.STATUS: utils.OK,
		utils.ERROR:  false,
	})

}

// Get detailed list of NSFW inscriptions godoc
// @Summary      List of Inscriptions in the wallet waiting to be approved
// @Description  List of Inscriptions in the wallet waiting to be approved
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success      200  {object}  models.NSFWInscriptionsResponse
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /inscriptions/nsfw [get]
func getNsfwInscription(c *fiber.Ctx) error {
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
	var res []models.NSFWTable
	var errDB error
	if pageSize == 0 && page == 0 {
		utils.ReportMessage(fmt.Sprintf("Get NSFW INS -> Offset: all, Limit: all"))
		res, errDB = db.ReadArrayStruct[models.NSFWTable]("SELECT * FROM NSFW_ORD WHERE approved = 0")
	} else {
		utils.ReportMessage(fmt.Sprintf("Get TX -> Offset: %d, Page: %d", page, req.PageSize))
		res, errDB = db.ReadArrayStruct[models.NSFWTable](`SELECT * FROM NSFW_ORD
                                                                WHERE oid NOT IN ( SELECT oid FROM TRANSACTIONS_ORD
                                                                ORDER BY id LIMIT ? ) AND approved = 0
                                                                ORDER BY id LIMIT ?`, page, pageSize)
	}
	if errDB != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	final := make([]models.TxTable, 0)
	for _, ins := range res {
		file := fmt.Sprintf("./data_final/%s.%s", ins.OrdID[:8], utils.InlineIF[string](strings.Split(ins.FileFormat, "/")[0] == "image", "webp", "txt"))
		b64, err := utils.ReadFileAsBase64(file)
		if err != nil {
			utils.WrapErrorLog(err.Error())
			continue
		}
		val := models.TxTable{
			ID:          ins.ID,
			OrdID:       ins.OrdID,
			TxID:        ins.TxID,
			FileFormat:  utils.InlineIF[string](strings.Split(ins.FileFormat, "/")[0] == "image", "image/webpwebp", "text/plain;charset=utf-8"),
			BcAddress:   ins.BcAddress,
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

// Estimate inscription cost godoc
// @Summary      Estimate inscription cost !!!Don't use this method!!!
// @Description  Estimate inscription cost !!!Don't use this method!!!
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param 		 data body models.EstimateRequest true "Image URL from hosting service and number of blocks"
// @Success      200  {object}  models.Inscribe
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /estimate [post]
func estimate(c *fiber.Ctx) error {
	var req models.EstimateRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ReportError(c, "Cannot parse body", http.StatusBadRequest)
	}
	if req.NumberOfBlocks < 1 {
		return utils.ReportError(c, "Number of blocks must be greater than 0", http.StatusBadRequest)
	}
	// Random number generator
	rand.NewSource(time.Now().UnixNano())
	randName := rand.Int()

	response, err := http.Get(req.ImageURL)
	if err != nil {
		return err
	}
	//if strings.Contains(typeFile, "image") {
	//    return utils.ReportError(c, "File must be picture", http.StatusBadRequest)
	//}
	fileName := fmt.Sprintf("%d.%s", randName, "webp")
	defer response.Body.Close()

	// Create an empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	// Make sure to delete the file after the function finishes
	defer os.Remove(fileName)
	utils.ReportMessage("File saved")

	dm := services.GetDaemon()
	sv, err := coind.WrapDaemon(dm, 1, "estimatesmartfee", req.NumberOfBlocks, "economical")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}
	var fRate models.FeeRate
	err = json.Unmarshal(sv, &fRate)
	if err != nil {
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}
	utils.ReportMessage(fmt.Sprintf("Fee rate: %f", fRate.Feerate))

	feeRate := int(fRate.Feerate / 1024 * 100000000)
	utils.ReportMessage(fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300  inscribe --dry-run --fee-rate %d %s", feeRate, fileName))
	s, err := cmd.CallJSON[models.Inscribe]("bash", "-c", fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300  inscribe --dry-run --fee-rate %d --file %s", feeRate, fileName))
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(s)
}

// Send inscription godoc
// @Summary      Test picture for NSFW content
// @Description  Test picture for NSFW content
// @Tags         NSFW
// @Accept       json
// @Produce      json
// @Param 		 data body models.TestPicReq true "File in base64 and filename"
// @Success      200  {object}  models.TestPicResponse
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /nsfw [post]
func testPic(c *fiber.Ctx) error {
	var req models.TestPicReq
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ReportError(c, "Cannot parse body", http.StatusBadRequest)
	}
	tx := &grpcModels.NSFWRequest{
		Base64:   req.Base64,
		Filename: req.Filename,
	}
	res, err := grpcClient.DetectNSFW(tx)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	return c.JSON(&fiber.Map{
		"nsfwPicture": res.NsfwPicture,
		"nsfwText":    res.NsfwText,
		utils.STATUS:  utils.OK,
		utils.ERROR:   false,
	})
}

// Get Raw transaction godoc
// @Summary      Get Raw transaction from BTC code
// @Description  Get Raw transaction from BTC code
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param tx query string true "Transaction ID"
// @Success      200  {object}  models.RawTransaction
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /transaction/raw [get]
func getRawTx(c *fiber.Ctx) error {
	tx := c.Query("tx", "")
	if tx == "" {
		return utils.ReportError(c, "Missing tx", http.StatusBadRequest)
	}
	utils.ReportMessage(fmt.Sprintf("Getting raw transaction %s", tx))
	dm := services.GetDaemon()
	sv, err := coind.WrapDaemon(dm, 1, "getrawtransaction", tx, true)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, "Cannot get raw transaction", http.StatusInternalServerError)
	}
	var trans models.RawTransaction
	err = json.Unmarshal(sv, &trans)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, "Cannot get raw transaction", http.StatusInternalServerError)
	}
	//rs := models.RawTxResponse{
	//	RawTx:    trans,
	//	HasError: false,
	//	Status:   utils.OK,
	//}
	return c.JSON(trans)
}

// Get Base64 image from Inscription ID godoc
// @Summary      Get Base64 image from Inscription ID
// @Description  Get Base64 image from Inscription ID
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param idInscription query string true "Inscription ID"
// @Success      200  {object}  models.InscriptionPicResponse
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /inscription/image [get]
func getImage(c *fiber.Ctx) error {
	tx := c.Query("idInscription", "")
	if tx == "" {
		return utils.ReportError(c, "Missing Inscription ID", http.StatusBadRequest)
	}
	utils.ReportMessage(fmt.Sprintf("Getting image for inscription %s", tx))
	file := "./data_final/" + tx[:8] + ".webp"
	b64, err := utils.ReadFileAsBase64(file)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return err
	}
	re := models.InscriptionPicResponse{
		HasError: false,
		Base64:   b64,
		Status:   utils.OK,
	}
	return c.JSON(re)
}

// Get fee rate for transaction godoc
// @Summary      Get fee rate for transaction
// @Description  Get fee rate for transaction
// @Tags         Fees
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.FeeResponse
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /feerate [get]
func getFeeRate(c *fiber.Ctx) error {
	dm := services.GetDaemon()
	sv, err := coind.WrapDaemon(dm, 1, "estimatesmartfee", 5, "economical")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}
	utils.ReportMessage(fmt.Sprintf("Fee rate %s sats", sv))
	var fRate models.FeeRate
	err = json.Unmarshal(sv, &fRate)
	if err != nil {
		return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	}

	feeRate := int(fRate.Feerate / 1024 * 100000000)
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"feerate":    feeRate,
		utils.STATUS: utils.OK,
		utils.ERROR:  false,
	})
}

// Get new address godoc
// @Summary      Get new BTC Address
// @Description  Get new BTC Address
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
	addr, err := cmd.CallJSONNonLock[Address]("bash", "-c", "/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300  wallet receive")
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	utils.ReportMessage(fmt.Sprintf("Get addr: %s", addr.Address))
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
	utils.ReportMessage(fmt.Sprintf("Send inscription id: %s", req.InscriptionID))
	if req.Address == "" {
		return utils.ReportError(c, "Address is empty", http.StatusBadRequest)
	}
	if req.InscriptionID == "" {
		return utils.ReportError(c, "Inscription id is empty", http.StatusBadRequest)
	}
	if req.FeeRate == 0 {
		return utils.ReportError(c, "FeeRate estimation error", http.StatusBadRequest)
	}

	s, err := cmd.CallString("bash", "-c", fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300  wallet send --fee-rate %d %s %s", req.FeeRate, req.Address, req.InscriptionID))
	if err != nil {
		utils.WrapErrorLog(err.Error())
		_, errWeb := utils.POSTRequest[models.ErrorHTTP]("https://rocket.art/api/api/v1/transfer/error", &fiber.Map{
			"error": err.Error(),
			"id":    req.ID,
		})
		if errWeb != nil {
			utils.WrapErrorLog(errWeb.Error())
		}
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)

	}
	utils.ReportMessage(fmt.Sprintf("Ordinal sent: %v", s))
	_, errWeb := utils.POSTRequest[models.ErrorHTTP]("https://rocket.art/api/api/v1/transfer/confirm", &fiber.Map{
		"tx_id": s,
		"id":    req.ID,
	})
	if errWeb != nil {
		utils.WrapErrorLog(errWeb.Error())
	}

	go func() {
		fileFormat := db.ReadValueEmpty[string]("SELECT file_format FROM NSFW_ORD WHERE ord_id = ?", req.InscriptionID)
		_, err = db.InsertSQl("DELETE FROM TRANSACTIONS_ORD WHERE ord_id = ?", req.InscriptionID)
		if err != nil {
			utils.WrapErrorLog("Can't delete inscription from DB")
			return
		}
		file := fmt.Sprintf("./data_final/%s.%s", req.InscriptionID[:8], utils.InlineIF[string](strings.Split(fileFormat, "/")[0] == "image", "webp", "txt"))
		err = os.Remove(file)
		if err != nil {
			utils.WrapErrorLog("Error deleting inscription from file system")
			return
		}
	}()

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
	utils.ReportMessage(fmt.Sprintf("Mint inscription id"))
	if req.Format == "" {
		return utils.ReportError(c, "Format is empty", http.StatusBadRequest)
	}
	if req.Base64 == "" {
		return utils.ReportError(c, "Base64 is empty", http.StatusBadRequest)
	}

	if req.FeeRate == 0 {
		return utils.ReportError(c, "Fee rate is empty", http.StatusBadRequest)
	}

	//dm := services.GetDaemon()
	//sv, err := coind.WrapDaemon(dm, 1, "estimatesmartfee", 5, "economical")
	//if err != nil {
	//	utils.WrapErrorLog(err.Error())
	//	return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	//}
	//var fRate models.FeeRate
	//err = json.Unmarshal(sv, &fRate)
	//if err != nil {
	//	return utils.ReportError(c, "Cannot estimate fee rate", http.StatusInternalServerError)
	//}
	//
	//feeRate := int(fRate.Feerate / 1024 * 100000000)

	byteArray, err := utils.DecodePayload([]byte(req.Base64))
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}
	//fileType := strings.Split(req.Format, "/")[0]
	fileName := fmt.Sprintf("temp.%s", req.Format)

	err = os.WriteFile(fileName, byteArray, 0644)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	tx := &grpcModels.NSFWRequest{
		Base64:   req.Base64,
		Filename: fmt.Sprintf("pic.%s", req.Format),
	}
	res, err := grpcClient.DetectNSFW(tx)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, "Cannot check if NSFW image", http.StatusInternalServerError)
	}

	if res.NsfwPicture {
		err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", fileName))).Run()
		if err != nil {
			utils.WrapErrorLog("Can't delete file in data")
		}
		return utils.ReportError(c, "NSFW image", http.StatusConflict)
	}

	if res.NsfwText {
		err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", fileName))).Run()
		if err != nil {
			utils.WrapErrorLog("Can't delete file in data")
		}
		return utils.ReportError(c, "NSFW Text in the image", http.StatusConflict)
	}
	utils.ReportMessage(fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300  inscribe --fee-rate %d --file %s", req.FeeRate, fileName))
	s, err := cmd.CallJSON[models.Inscribe]("bash", "-c", fmt.Sprintf("/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300  inscribe --fee-rate %d --file %s", req.FeeRate, fileName))
	if err != nil {
		if strings.Contains(strings.ToUpper(err.Error()), "EXIT STATUS 1") {
			go func() {
				err := exec.Command("bash", "-c", "systemctl --user restart bitcoin").Run()
				if err != nil {
					utils.WrapErrorLog(err.Error())
				}
				time.Sleep(120 * time.Second)
				err = exec.Command("bash", "-c", "/home/dfwplay/bitcoin-cli loadwallet ord").Run()
				if err != nil {
					utils.WrapErrorLog(err.Error())
				}
			}()
		}
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
		utils.ReportMessage(fmt.Sprintf("Get TX -> Offset: all, Limit: all"))
		res, errDB = db.ReadArrayStruct[models.TxTable]("SELECT * FROM TRANSACTIONS_ORD ORDER BY CAST(id AS INTEGER) DESC")
	} else {
		utils.ReportMessage(fmt.Sprintf("Get TX -> Offset: %d, Page: %d", page, req.PageSize))
		res, errDB = db.ReadArrayStruct[models.TxTable](`SELECT * FROM TRANSACTIONS_ORD
WHERE oid NOT IN ( SELECT oid FROM TRANSACTIONS_ORD
                   ORDER BY CAST(id AS INTEGER) DESC LIMIT ? )
 LIMIT ?`, page, pageSize)
	}
	if errDB != nil {
		utils.WrapErrorLog(err.Error())
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	final := make([]models.TxTable, 0)
	var wg sync.WaitGroup
	wg.Add(len(res))
	for _, ins := range res {
		go func(ins models.TxTable) {
			defer wg.Done()
			if len(ins.OrdID) == 0 {
				return
			}
			file := fmt.Sprintf("./data_final/%s.%s", ins.OrdID[:8], utils.InlineIF[string](strings.Split(ins.FileFormat, "/")[0] == "image", "webp", "txt"))
			b64, err := utils.ReadFileAsBase64(file)
			if err != nil {
				utils.WrapErrorLog(err.Error())
				return
			}
			val := models.TxTable{
				ID:          ins.ID,
				OrdID:       ins.OrdID,
				TxID:        ins.TxID,
				FileFormat:  utils.InlineIF[string](strings.Split(ins.FileFormat, "/")[0] == "image", "image/webp", "text/plain;charset=utf-8"),
				BcAddress:   ins.BcAddress,
				Link:        ins.Link,
				ContentLink: ins.ContentLink,
				Base64:      b64,
			}
			final = append(final, val)
		}(ins)
	}
	wg.Wait()

	sort.Slice(final, func(i, j int) bool {
		return final[i].ID > final[j].ID
	})

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
	utils.ReportMessage(fmt.Sprintf("Get TX -> Offset: %d, Page: %d", page, req.PageSize))

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

// EstimateFeeQuality Estimate inscription cost based on quality godoc
// @Summary      Estimate inscription cost based on quality
// @Description  Estimate inscription cost based on quality
// @Tags         Inscriptions
// @Accept       json
// @Produce      json
// @Param 		 data body models.EstimateQualityRequest true "Image URL from hosting service and number of blocks"
// @Success      200  {object}  models.EstimateQualityResponse
// @Failure      400  {object}  models.ErrorHTTP
// @Failure      409  {object}  models.ErrorHTTP
// @Failure      500  {object}  models.ErrorHTTP
// @Router       /fee/quality/estimate [post]
func EstimateFeeQuality(c *fiber.Ctx) error {
	var req models.EstimateQualityRequest
	err := c.BodyParser(&req)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusBadRequest)
	}
	url := req.UrlPic      // url of the image
	quality := req.Quality // Quality to be passed to cwebp

	// Ensure temp folder exists
	err = os.MkdirAll("temp", os.ModePerm)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}

	// Create a temp file in the temp folder
	tempFile, err := os.CreateTemp("temp", "download-*.webp")
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	defer os.Remove(tempFile.Name()) // clean up

	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	tempFile.Close()

	// Change filename for webp file
	outputFile := tempFile.Name() + "-processed"

	// Convert image to webp
	x := exec.Command("cwebp", "-q", quality, "-resize", "500", "0", tempFile.Name(), "-o", outputFile)
	err = x.Run()
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}
	defer os.Remove(outputFile)

	// Read the webp image file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		return utils.ReportError(c, err.Error(), http.StatusInternalServerError)
	}

	// Convert to base64
	str := base64.StdEncoding.EncodeToString(data)

	// Get the size of the file in bytes. Use FileInfo from os package.
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	// Size of file in bytes
	sizeInBytes := fileInfo.Size()

	//calculate ordinal fee
	feeRate := req.FeeRate
	cin := 2.0
	cout := 3.0
	precalc := 144.5 + (cin * 57.5) + (cout * 43.0) + float64(sizeInBytes)/3.977*float64(feeRate) + 10000
	btcAmount := precalc / 100_000_000.0 * 1.05

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"size":       sizeInBytes,
		"base64":     str,
		"btc_amount": btcAmount,
	})
}
