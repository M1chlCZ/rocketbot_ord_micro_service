package db

import (
	"api/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type DB struct {
	*sqlx.DB
}

var MysqlDB *DB

func InitMySQL() {
	db, errDB := sqlx.Open("mysql", utils.GetENV("DB_CONN"))
	if errDB != nil {
		log.Fatal(errDB)
		return
	}
	MysqlDB = &DB{db}
}
func MySQLReadSql(SQL string, params ...interface{}) (*sqlx.Rows, error) {
	results, errRow := MysqlDB.Queryx(SQL, params...)
	if errRow != nil {
		fmt.Println(errRow.Error())
		return nil, errRow
	} else {
		return results, nil
	}
}

func MySQLReadValue[T any](SQL string, params ...interface{}) (T, error) {
	d := make(chan T, 1)
	e := make(chan error, 1)
	go func(data chan T, err chan error) {
		var an T
		errDB := MysqlDB.QueryRow(SQL, params...).Scan(&an)
		if errDB != nil {
			err <- errDB
		} else {
			data <- an
		}
	}(d, e)
	select {
	case data := <-d:
		close(d)
		close(e)
		return data, nil
	case err := <-e:
		close(d)
		close(e)
		return getZero[T](), err
	}
}

func MySQLReadValueEmpty[T any](SQL string, params ...interface{}) T {
	d := make(chan T, 1)
	e := make(chan error, 1)
	go func(data chan T, err chan error) {
		var an T
		errDB := MysqlDB.QueryRow(SQL, params...).Scan(&an)
		if errDB != nil {
			err <- errDB
		} else {
			data <- an
		}
	}(d, e)
	select {
	case data := <-d:
		close(d)
		close(e)
		return data
	case err := <-e:
		close(d)
		close(e)
		log.Println(err)
		return getZero[T]()
	}
}

func MySQLReadStruct[T any](SQL string, params ...interface{}) (T, error) {
	d := make(chan T, 1)
	e := make(chan error, 1)
	go func(data chan T, err chan error) {
		rows, errDB := MysqlDB.Queryx(SQL, params...)
		if errDB != nil {
			_ = rows.Close()
			err <- errDB
		} else {
			var s T
			s, errDB := ParseStruct[T](rows)
			if errDB != nil {
				_ = rows.Close()
				err <- errDB
			}
			_ = rows.Close()
			data <- s
		}
	}(d, e)
	select {
	case data := <-d:
		close(d)
		close(e)
		return data, nil
	case err := <-e:
		close(d)
		close(e)
		return getZero[T](), err
	}
}

func MySQLReadStructEmpty[T any](SQL string, params ...interface{}) T {
	d := make(chan T, 1)
	go func(data chan T) {
		rows, err := MysqlDB.Queryx(SQL, params...)
		if err != nil {
			utils.WrapErrorLog(err.Error())
			i := getZero[T]()
			_ = rows.Close()
			data <- i
			//return i
		} else {
			var s T
			s, err := ParseStruct[T](rows)
			if err != nil {
				utils.WrapErrorLog(err.Error())
				_ = rows.Close()
				data <- getZero[T]()
				//return getZero[T]()
			}
			_ = rows.Close()
			data <- s
		}
	}(d)
	select {
	case data := <-d:
		close(d)
		return data
	}
}

func MySQLReadArrayStruct[T any](SQL string, params ...interface{}) ([]T, error) {
	d := make(chan []T, 1)
	e := make(chan error, 1)
	go func(data chan []T, err chan error) {
		rows, errDB := ReadSql(SQL, params...)
		if errDB != nil {
			//utils.WrapErrorLog(err.Error())
			//i := getZeroArray[T]()
			//data <- i
			err <- errDB
		} else {
			s := ParseArrayStruct[T](rows)
			if errDB != nil {
				_ = rows.Close()
				err <- errDB
			}
			_ = rows.Close()
			data <- s
		}
	}(d, e)
	select {
	case data := <-d:
		close(d)
		close(e)
		return data, nil
	case err := <-e:
		close(d)
		close(e)
		return getZeroArray[T](), err
	}
}

func MySQLReadArray[T any](SQL string, params ...interface{}) ([]T, error) {
	d := make(chan []T, 1)
	e := make(chan error, 1)
	go func(data chan []T, err chan error) {
		i := make([]T, 0)
		rows, errDB := MysqlDB.Queryx(SQL, params...)
		if errDB != nil {
			utils.WrapErrorLog(errDB.Error())
			//data <- i
			err <- errDB
		} else {
			for rows.Next() {
				var s T
				if errDB := rows.Scan(&s); errDB != nil {
					//data <- i
					err <- errDB
				} else {
					i = append(i, s)
				}
			}
			_ = rows.Close()
			data <- i
		}
	}(d, e)
	select {
	case data := <-d:
		close(d)
		close(e)
		return data, nil
	case err := <-e:
		close(d)
		close(e)
		return getZeroArray[T](), err
	}
}

func MySQLInsert(SQL string, params ...interface{}) (int64, error) {
	query, errStmt := MysqlDB.Exec(SQL, params...)
	if errStmt != nil {
		//fmt.Printf("Can't Insert shit")
		return 0, errStmt
	}
	id, errLastID := query.LastInsertId()
	if errLastID != nil {
		return 0, errLastID
	}
	return id, nil
}
