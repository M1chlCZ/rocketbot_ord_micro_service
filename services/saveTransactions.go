package services

import (
	"api/cmd"
	"api/coind"
	"api/db"
	"api/models"
	"api/utils"
	"database/sql"
	"encoding/json"
	"strings"
)

func SaveListTransaction() {
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
