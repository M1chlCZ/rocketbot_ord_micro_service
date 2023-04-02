package db

import (
	"api/utils"
	"database/sql"
	"fmt"
	_ "github.com/mutecomm/go-sqlcipher"
	"net/url"
	"os"
	"strconv"
)

type DBClient struct {
	client *sql.DB
}

var colorReset = "\033[0m"

const dbName string = "./.datb"

const dbVersion int = 1

var dbClient DBClient

func InitDB() (*DBClient, error) {
	if dbClient.client != nil {
		return &dbClient, nil
	}
	utils.ReportMessage("DB opening")
	key := url.QueryEscape("kGMbPd3BrGJ6Htd")
	dbname := fmt.Sprintf("%s?_pragma_key=%s&_pragma_cipher_page_size=4096", dbName, key)

	exists := false
	if _, err := os.Stat(dbName); err != nil {
		exists = false
	} else {
		exists = true
	}

	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		err := os.Remove(dbName)
		if err != nil {
			utils.WrapErrorLog(err.Error())
			return nil, err
		}
		utils.WrapErrorLog(err.Error())
		return nil, err
	}

	if !exists {
		_ = ExecQuery(db, fmt.Sprintf("PRAGMA user_version = %d", dbVersion))
	}
	initTables(db)
	dbClient = DBClient{client: db}
	return &dbClient, nil
}

func initTables(db *sql.DB) {
	createTokenTable := `CREATE TABLE IF NOT EXISTS TOKEN_TABLE (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"token" TEXT
		
	  );`

	createJWTTable := `CREATE TABLE IF NOT EXISTS JWT_TABLE (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"JWT" TEXT,
		"refreshToken" TEXT
	  );`

	createDaemonTable := `CREATE TABLE IF NOT EXISTS DAEMON_TABLE (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"wallet_usr" TEXT NOT NULL,
		"wallet_pass" TEXT NOT NULL,
		"wallet_port" INTEGER NOT NULL,
		"folder" TEXT NOT NULL,
		"node_id" INTEGER NOT NULL,
		"coin_id" INTEGER NOT NULL,
		"conf" TEXT NOT NULL,
		"mn_port" INT NOT NULL,
		"ip" TEXT NOT NULL,
		"wallet_passphrase" TEXT
	  );`

	createStakingDeamonTable := `CREATE TABLE IF NOT EXISTS STAKING_DAEMON_TABLE (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"wallet_usr" TEXT NOT NULL,
		"wallet_pass" TEXT NOT NULL,
		"wallet_port" INTEGER NOT NULL,
		"coin_id" INTEGER NOT NULL,
		"ip" TEXT NOT NULL,
		"wallet_passphrase" TEXT
	  );`

	err := ExecQuery(db, createJWTTable)
	err = ExecQuery(db, createTokenTable)
	err = ExecQuery(db, createDaemonTable)
	err = ExecQuery(db, createStakingDeamonTable)

	err, i := GetVersion(db)
	if err != nil {
		return
	}

	switch i {
	case 0, 1:
		err = ExecQuery(db, "ALTER TABLE DAEMON_TABLE ADD COLUMN conf TEXT NOT NULL DEFAULT ('')")
		err = ExecQuery(db, "ALTER TABLE DAEMON_TABLE ADD COLUMN ip TEXT NOT NULL DEFAULT ('')")
		break
	case 2:
		err = ExecQuery(db, "ALTER TABLE DAEMON_TABLE ADD COLUMN mn_port INT NOT NULL DEFAULT 0")
		break
	case 3:
		err = ExecQuery(db, "ALTER TABLE DAEMON_TABLE ADD COLUMN wallet_passphrase TEXT")
		break
	case 4:
		err = ExecQuery(db, "ALTER TABLE JWT_TABLE ADD COLUMN refreshToken TEXT")
		break
	case 5:
		err = ExecQuery(db, "ALTER TABLE STAKING_DAEMON_TABLE ADD COLUMN ip TEXT DEFAULT ('127.0.0.1')")
		break
	default:
		break
	}

	err = ExecQuery(db, fmt.Sprintf("PRAGMA user_version = %d", dbVersion))

	if err != nil {
		fmt.Printf("Error while creating table token")
		fmt.Println(err.Error())
		return
	}

}

func ExecQuery(db *sql.DB, sql string) error {
	statementJWT, err := db.Prepare(sql) // Prepare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statementJWT.Exec()
	if err != nil {
		fmt.Printf("Error while creating table jwt")
		fmt.Println(err.Error())
		return err
	}
	_ = statementJWT.Close()
	return nil
}

func GetVersion(db *sql.DB) (error, int) {
	insertStudentSQL := `PRAGMA user_version`
	rows := db.QueryRow(insertStudentSQL)
	var Ver string
	err := rows.Scan(&Ver)

	if err != nil {
		return err, 0
		//log.Fatalln(err.Error())
	}

	atoi, err := strconv.Atoi(Ver)
	if err != nil {
		return err, 0
	}
	return nil, atoi
}
