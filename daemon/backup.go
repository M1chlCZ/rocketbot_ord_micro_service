package daemon

import (
	"api/coind"
	"api/services"
	"api/utils"
	"fmt"
	"os/exec"
)

func BackupBitcoinWallet() {
	utils.ReportMessage("Backing up bitcoin wallet")
	daemon := services.GetDaemon()

	_, err := coind.WrapDaemon(daemon, 1, "backupwallet", "api/backup.db")
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}

	err = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("zip %s %s", "backup.zip", "backup.db"))).Run()
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	err = utils.SendBackupToServer(fmt.Sprintf("%s/api/%s", utils.GetHomeDir(), "backup.zip"))
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return
	}
	utils.ReportMessage("Backup sent to server")

	_ = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", "backup.db"))).Run()
	_ = exec.Command("bash", "-c", fmt.Sprintf(fmt.Sprintf("rm %s", "backup.zip"))).Run()

}
