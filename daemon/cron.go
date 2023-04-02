package daemon

import (
	"api/cmd"
	"api/utils"
	"github.com/jasonlvhit/gocron"
)

func StartCron() {
	err := gocron.Every(10).Minutes().Do(index)

	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	//utils.ReportMessage("< - Cron service started - >")
	<-gocron.Start()

}

func index() {
	utils.ReportMessage("Indexing")
	a, err := cmd.CallString("bash", "-c", "/home/dfwplay/bin/ord --cookie-file ~/.bitcoin/.cookie --rpc-url 127.0.0.1:12300/wallet/ord --wallet ord wallet inscriptions")
	if err != nil {
		utils.WrapErrorLog(err.Error())
	}
	utils.ReportMessage(string(a))
}
