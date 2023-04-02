package main

import (
	"api/cmd"
	"api/coind"
	"api/db"
	"api/models"
	"api/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {

	err := db.InitDB()
	if err != nil {
		utils.WrapErrorLog(fmt.Sprintf("Error opening db: %s", err.Error()))
		return
	}

	utils.ReportMessage("DB opened")

	s2, err := cmd.CallString("bash", "-c", "cat /home/dfwplay/.bitcoin/.cookie")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}

	s := strings.Split(s2, ":")
	utils.ReportMessage(s[1])
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
			utils.WrapErrorLog(err.Error())
			continue
		}

		//utils.ReportMessage(fmt.Sprintf("result: %s %s %s %s %s", ins.Inscription, txid[0], addr, ins.Location, ins.Explorer))

	}

	utils.ReportMessage("Inscriptions saved into db")

	time.Sleep(time.Second * 1)
}
