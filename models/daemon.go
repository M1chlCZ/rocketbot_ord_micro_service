package models

import (
	"database/sql"
	"fmt"
)

type DaemonCommon interface {
	GetDaemon() *Daemon
	IsRemote() bool
}

type Daemon struct {
	ID         int            `db:"id"`
	WalletUser string         `db:"wallet_usr"`
	WalletPass string         `db:"wallet_pass"`
	WalletPort int            `db:"wallet_port"`
	Wallet     string         `db:"wallet"`
	Folder     string         `db:"folder"`
	NodeID     int            `db:"node_id"`
	CoinID     int            `db:"coin_id"`
	Conf       string         `db:"conf"`
	IP         string         `db:"ip"`
	MnPort     int            `db:"mn_port"`
	PassPhrase sql.NullString `db:"wallet_passphrase"`
}

func (d *Daemon) ToString() string {
	return fmt.Sprintf("[%s, %s, %d, %s, %d, %d, %s, %s, %d]", d.WalletUser, d.WalletPass, d.WalletPort, d.Folder, d.NodeID, d.CoinID, d.Conf, d.IP, d.MnPort)
}
func (d *Daemon) GetDaemon() *Daemon {
	return d
}

func (d *Daemon) IsRemote() bool {
	return false
}

type BitcoinDaemon struct {
	ID         int            `db:"id"`
	WalletUser string         `db:"wallet_usr"`
	WalletPass string         `db:"wallet_pass"`
	WalletPort int            `db:"wallet_port"`
	Wallet     string         `db:"wallet"`
	CoinID     int            `db:"coin_id"`
	PassPhrase sql.NullString `db:"wallet_passphrase"`
	IP         string         `db:"ip"`
}

func (d *BitcoinDaemon) ToString() string {
	return fmt.Sprintf(`[%s, %s, %d, %d, %s]`, d.WalletUser, d.WalletPass, d.WalletPort, d.CoinID, d.PassPhrase.String)
}

func (d *BitcoinDaemon) GetDaemon() *Daemon {
	daemon := &Daemon{
		WalletUser: d.WalletUser,
		WalletPass: d.WalletPass,
		WalletPort: d.WalletPort,
		Wallet:     d.Wallet,
		CoinID:     d.CoinID,
		PassPhrase: d.PassPhrase,
		IP:         d.IP,
	}
	return daemon
}

func (d *BitcoinDaemon) IsRemote() bool {
	return d.IP != "" && d.IP != "127.0.0.1"
}
