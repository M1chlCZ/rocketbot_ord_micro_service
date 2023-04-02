package coind

import (
	"api/models"
	"api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

func WrapDaemon(daemon models.DaemonCommon, maxTries int, method string, params ...interface{}) ([]byte, error) {
	wg := new(sync.WaitGroup)
	c := make(chan []byte, 1)
	e := make(chan error, 1)
	var res []byte
	var err error
	go callDaemon(c, e, wg, daemon, maxTries, method, params)
	wg.Wait()
	select {
	case data := <-c:
		res = data
	case er := <-e:
		utils.WrapErrorLog(er.Error())
		err = er
	}
	close(c)
	close(e)

	if err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

func callDaemon(c chan []byte, e chan error, wg *sync.WaitGroup, daemonCommon models.DaemonCommon, triesMax int, command string, params any) {
	wg.Add(1)
	defer wg.Done()
	var client *Coind
	var errClient error
	tries := 0
	daemon := daemonCommon.GetDaemon()
	ip := "127.0.0.1"
	if daemonCommon.IsRemote() != false {
		ip = daemon.IP
	}
	//utils.ReportMessage(fmt.Sprintf("Calling %s %s", daemon.Folder, command))
	for {
		if tries > 3 {
			utils.ReportMessage(fmt.Sprintf("Try %d of %d. CMD: %s Daemon: %s", tries, triesMax, command, daemon.Folder))
		}
		if tries >= triesMax {
			if errClient != nil {
				e <- errClient
			} else {
				e <- errors.New("error getting RPC data")
			}
			return
		}

		tries++
		if client == nil {
			client, errClient = New(ip, daemon.WalletPort, daemon.Wallet, daemon.WalletUser, daemon.WalletPass, false)
			if errClient != nil {
				errClient = errors.New(fmt.Sprintf("error connecting to %s wallet daemon: %s", daemon.Folder, errClient.Error()))
				utils.ReportMessage(errClient.Error())
				time.Sleep(30 * time.Second)
				continue
			}
		}
		p, errClient := client.Call(command, params)

		if errClient != nil {
			var er models.ErrorWallet
			err := json.Unmarshal([]byte(errClient.Error()), &er)
			if err != nil {
				utils.ReportWarning(err.Error())
				time.Sleep(30 * time.Second)
				continue
			} else {
				utils.ReportMessage("JSON error: " + er.Error.Message)
				if er.Error.Code == -28 {
					utils.ReportWarning(fmt.Sprintf("|/> wallet daemon %s is still syncing <\\|", daemon.Folder))
					time.Sleep(180 * time.Second)
					continue
				}
				if er.Error.Code == -1 {
					e <- errors.New(er.Error.Message)
					return
				}
				time.Sleep(30 * time.Second)
				continue
			}

		}
		if string(p) != "null" {
			if len(p) != 0 {
				c <- p
				return
			}

			//utils.ReportMessage("Error, trying again")
			time.Sleep(35 * time.Second)
		}
	}
}
