package models

type GetWalletInfo struct {
	Walletname            string  `json:"walletname"`
	Walletversion         int     `json:"walletversion"`
	Format                string  `json:"format"`
	Balance               float64 `json:"balance"`
	UnconfirmedBalance    float64 `json:"unconfirmed_balance"`
	ImmatureBalance       float64 `json:"immature_balance"`
	Txcount               int     `json:"txcount"`
	Keypoolsize           int     `json:"keypoolsize"`
	KeypoolsizeHdInternal int     `json:"keypoolsize_hd_internal"`
	Paytxfee              float64 `json:"paytxfee"`
	PrivateKeysEnabled    bool    `json:"private_keys_enabled"`
	AvoidReuse            bool    `json:"avoid_reuse"`
	Scanning              bool    `json:"scanning"`
	Descriptors           bool    `json:"descriptors"`
	ExternalSigner        bool    `json:"external_signer"`
}

type ListUnspent struct {
	Txid          string   `json:"txid"`
	Vout          int      `json:"vout"`
	Address       string   `json:"address"`
	ScriptPubKey  string   `json:"scriptPubKey"`
	Amount        float64  `json:"amount"`
	Confirmations int      `json:"confirmations"`
	Spendable     bool     `json:"spendable"`
	Solvable      bool     `json:"solvable"`
	Desc          string   `json:"desc"`
	ParentDescs   []string `json:"parent_descs"`
	Safe          bool     `json:"safe"`
	Label         string   `json:"label,omitempty"`
}

type GetTransaction struct {
	Amount            float64       `json:"amount"`
	Fee               float64       `json:"fee"`
	Confirmations     int           `json:"confirmations"`
	Blockhash         string        `json:"blockhash"`
	Blockheight       int           `json:"blockheight"`
	Blockindex        int           `json:"blockindex"`
	Blocktime         int           `json:"blocktime"`
	Txid              string        `json:"txid"`
	Wtxid             string        `json:"wtxid"`
	Walletconflicts   []interface{} `json:"walletconflicts"`
	Time              int           `json:"time"`
	Timereceived      int           `json:"timereceived"`
	Bip125Replaceable string        `json:"bip125-replaceable"`
	Details           []interface{} `json:"details"`
	Hex               string        `json:"hex"`
	Decoded           Decoded       `json:"decoded"`
}
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}
type Vin struct {
	Txid        string    `json:"txid"`
	Vout        int       `json:"vout"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Txinwitness []string  `json:"txinwitness"`
	Sequence    int64     `json:"sequence"`
}
type ScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}
type Decoded struct {
	Txid     string `json:"txid"`
	Hash     string `json:"hash"`
	Version  int    `json:"version"`
	Size     int    `json:"size"`
	Vsize    int    `json:"vsize"`
	Weight   int    `json:"weight"`
	Locktime int    `json:"locktime"`
	Vin      []Vin  `json:"vin"`
	Vout     []Vout `json:"vout"`
}

type ListTransactions []struct {
	Address           string        `json:"address"`
	Category          string        `json:"category"`
	Amount            float64       `json:"amount"`
	Vout              int           `json:"vout"`
	Fee               float64       `json:"fee,omitempty"`
	Confirmations     int           `json:"confirmations"`
	Blockhash         string        `json:"blockhash,omitempty"`
	Blockheight       int           `json:"blockheight,omitempty"`
	Blockindex        int           `json:"blockindex,omitempty"`
	Blocktime         int           `json:"blocktime,omitempty"`
	Txid              string        `json:"txid"`
	Wtxid             string        `json:"wtxid"`
	Walletconflicts   []interface{} `json:"walletconflicts"`
	Time              int           `json:"time"`
	Timereceived      int           `json:"timereceived"`
	Bip125Replaceable string        `json:"bip125-replaceable"`
	Abandoned         bool          `json:"abandoned,omitempty"`
	ParentDescs       []string      `json:"parent_descs,omitempty"`
	Label             string        `json:"label,omitempty"`
	Trusted           bool          `json:"trusted,omitempty"`
}

type ListTransactionsDB struct {
	ID                int           `json:"id" db:"id"`
	Address           string        `json:"address" db:"address"`
	Category          string        `json:"category" db:"category"`
	Amount            float64       `json:"amount" db:"amount"`
	Vout              int           `json:"vout" db:"vout"`
	Fee               float64       `json:"fee,omitempty" db:"fee"`
	Confirmations     int           `json:"confirmations" db:"confirmations"`
	Blockhash         string        `json:"blockhash,omitempty" db:"blockhash"`
	Blockheight       int           `json:"blockheight,omitempty" db:"blockheight"`
	Blockindex        int           `json:"blockindex,omitempty" db:"blockindex"`
	Blocktime         int           `json:"blocktime,omitempty" db:"blocktime"`
	Txid              string        `json:"txid" db:"txid"`
	Wtxid             string        `json:"wtxid" db:"wtxid"`
	Walletconflicts   []interface{} `json:"walletconflicts" db:"walletconflicts"`
	Time              int           `json:"time" db:"time"`
	Timereceived      int           `json:"timereceived" db:"timereceived"`
	Bip125Replaceable string        `json:"bip125-replaceable" db:"bip125_replaceable"`
	Abandoned         bool          `json:"abandoned,omitempty" db:"abandoned"`
	ParentDescs       []string      `json:"parent_descs,omitempty" db:"parent_descs"`
	Label             string        `json:"label,omitempty" db:"label"`
	Trusted           bool          `json:"trusted,omitempty" db:"trusted"`
}
