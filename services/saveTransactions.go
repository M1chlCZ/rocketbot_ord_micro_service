package services

import (
	"api/coind"
	"api/db"
	"api/models"
	"api/utils"
	"encoding/json"
	"strings"
)

func SaveListTransaction() {

	dm := GetDaemon()
	sv, err := coind.WrapDaemon(dm, 1, "listtransactions", "*", 9999999)
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
