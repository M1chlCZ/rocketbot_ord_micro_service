package services

import (
	"api/cmd"
	"api/models"
	"api/utils"
	"database/sql"
	"strings"
)

func GetDaemon() *models.BitcoinDaemon {
	s2, err := cmd.CallString("bash", "-c", "cat /home/dfwplay/.bitcoin/.cookie")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return nil
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
	return dm
}
