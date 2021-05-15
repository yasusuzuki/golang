package main

import (
	"database/sql" //SQLite
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type DBRecord map[string]interface{}

func queryForRows(db *sql.DB, sqlString string) ([]string, []DBRecord) {
	array := make([]DBRecord, 0, 5)
	var columns []string
	rows, err := db.Query(sqlString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, err = rows.Columns()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var row = make([]interface{}, len(columns))
		var rowp = make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			rowp[i] = &row[i]
		}

		rows.Scan(rowp...)

		rowMap := make(DBRecord)
		for i, col := range columns {
			switch row[i].(type) {
			case []byte:
				row[i] = string(row[i].([]byte))
				num, err := strconv.Atoi(row[i].(string))
				if err == nil {
					row[i] = num
				}
			}
			rowMap[col] = row[i]
		}
		array = append(array, rowMap)
	}

	return columns, array
}

func queryKeiyaku(db *sql.DB) {
	sql := "SELECT * FROM 契約" // WHERE type='table'"
	_, records := queryForRows(db, sql)
	fmt.Printf("%#v\n", records)
}

func showTables(db *sql.DB) {
	sql := "SELECT name FROM sqlite_master" // WHERE type='table'"
	_, records := queryForRows(db, sql)
	fmt.Printf("%#v\n", records)
}

var KEIYAKU_DATA []DBRecord

func initDB(db *sql.DB) {
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"ID": "D00000001", "client": "鈴木", "contractType": "A"})
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"ID": "D00000002", "client": "田中", "contractType": "A"})
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"ID": "D00000003", "client": "佐藤", "contractType": "B"})
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"ID": "D00000004", "client": "織田", "contractType": "B"})

	_, err := db.Exec(`CREATE TABLE "契約" ("ID" CHAR(12) PRIMARY KEY, "client" VARCHAR(64), "contractType" CHAR(4))`)
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare(`INSERT INTO "契約" ("ID", "client","contractType") VALUES (?, ?,?) `)
	if err != nil {
		panic(err)
	}
	for _, v := range KEIYAKU_DATA {
		if _, err = stmt.Exec(v["ID"], v["client"], v["contractType"]); err != nil {
			panic(err)
		}
	}

	stmt.Close()
}

func ConnectDB(dbfilename string) *sql.DB {
	db, err := sql.Open("sqlite3", dbfilename)
	if err != nil {
		panic(err)
	}
	return db
}

func createDB(dbfilename string) {
	//TODO
}

func main() {
	os.Remove("./data/test.db")
	createDB("aaa")
	db := ConnectDB("./data/test.db")
	initDB(db)
	showTables(db)
	queryKeiyaku(db)
}
