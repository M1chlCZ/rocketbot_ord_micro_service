package daemon

import (
	"api/cmd"
	"api/services"
	"api/utils"
	"github.com/jasonlvhit/gocron"
)

func StartCron() {
	err := gocron.Every(10).Minutes().Do(index)
	err = gocron.Every(1).Hour().Do(getInscriptions)
	err = gocron.Every(12).Hours().Do(BackupBitcoinWallet)

	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	//utils.ReportMessage("< - Cron service started - >")
	<-gocron.Start()

}

func index() {
	utils.ReportMessage("Indexing")
	_, err := cmd.CallString("bash", "-c", "/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300/wallet/ord --wallet ord wallet inscriptions")
	if err != nil {
		utils.WrapErrorLog(err.Error())
	}
}

func getInscriptions() {
	go services.GetInscriptions()
}
