package internal

import (
	"api/services"
	"api/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func SubmitTransaction(c *fiber.Ctx) error {
	//curl -X POST -H "txid:$1" -H "coinid:$coinID" -H "node_id:$nodeID" http://localhost:7500/submitTransaction
	txid := c.Get("txid")
	//
	//mp := &fiber.Map{
	//	"txid":   txid,
	//	"coinid": 0,
	//}

	utils.ReportMessage("Submitting transaction: " + txid)

	go services.GetInscriptions()
	go services.SaveListTransaction()

	//for {
	//	r, err := utils.POSTRequest[models.ErrorHTTP]("submitTransaction", mp)
	//	if err != nil {
	//		utils.WrapErrorLog(err.Error())
	//		time.Sleep(time.Second * 5)
	//		continue
	//	}
	//	if r.HasError != false {
	//		utils.WrapErrorLog(r.ErrorMessage)
	//		time.Sleep(time.Second * 5)
	//		continue
	//	}
	//	break
	//}
	////TODO THIS

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		utils.STATUS: utils.OK,
	})
}
