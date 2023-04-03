package models

type TxTable struct {
	ID          int    `db:"id" json:"id"`
	OrdID       string `db:"ord_id" json:"ord_id"`
	TxID        string `db:"tx_id" json:"tx_id"`
	BcAddress   string `db:"bc_address" json:"bc_address"`
	Link        string `db:"link" json:"link"`
	ContentLink string `db:"content_link" json:"content_link"`
	Base64      string `db:"-" json:"base64,omitempty"`
}
