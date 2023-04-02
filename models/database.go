package models

type TxTable struct {
	ID          int    `db:"id"`
	OrdID       string `db:"ord_id"`
	TxID        string `db:"tx_id"`
	BcAddress   string `db:"bc_address"`
	Link        string `db:"link"`
	ContentLink string `db:"content_link"`
}
