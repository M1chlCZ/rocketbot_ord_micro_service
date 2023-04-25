package services

import (
	"api/cmd"
	"api/coind"
	"api/db"
	"api/models"
	"api/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func GetInscriptions() {
	dm := GetDaemon()

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
		sv, err := coind.WrapDaemon(dm, 1, "gettransaction", txid[0], false, true)
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

		_, err = utils.DownloadImage(ins.Inscription)
		if err != nil {
			utils.WrapErrorLog(err.Error())
		}
		contentLink := fmt.Sprintf("https://ordinals.com/content/%s", ins.Inscription)
		addr := info.Decoded.Vout[vout].ScriptPubKey.Address
		_, err = db.InsertSQl(`INSERT INTO TRANSACTIONS_ORD (tx_id, ord_id, bc_address, link, content_link) 
									VALUES (?,?, ?, ?, ?)`, txid[0], ins.Inscription, addr, ins.Explorer, contentLink)

		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				continue
			}
			utils.WrapErrorLog(err.Error())
			continue
		}

	}
	ScanAndConvert()
	utils.ReportMessage("Inscriptions saved into db")
}
