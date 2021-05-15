// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/alexbrainman/odbc"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

/**
Access 2003 以前は、JET (Joint Engine Technology) データベースエンジン
Access 2007以降は、ACE (Access Connectivity Engine) データベースエンジン
https://qiita.com/yaju/items/7b0aa9e9f30005f60388
http://surferonwww.info/BlogEngine/post/2011/11/08/Development-of-application-which-uses-accdb-file-of-Access-2007.aspx


補足：　OLEとODBCの違い
https://mat0401.info/blog/dao-ado-odbc-oledb/
**/
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
	//TODO: やりかたがわからない
	//sql := "SELECT name FROM sqlite_master" // WHERE type='table'"
	//_, records := queryForRows(db, sql)
	//fmt.Printf("%#v\n", records)
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
	// ODBCをつかったDB接続
	// "C:\WINDOWS\system32\odbcad32.exe" で確認したドライバ名を"DRIVER"の欄にいれる必要あり
	db, err := sql.Open("odbc", "DRIVER={Microsoft Access Driver (*.mdb, *.accdb)};DBQ="+dbfilename+";")

	if err != nil {
		panic(err)
	}
	return db
}

/**
参考資料：
https://github.com/alexbrainman/odbc/blob/master/access_test.go
*/
func createDB(dbfilename string) {
	err := ole.CoInitialize(0)
	if err != nil {
		log.Fatal(err)
	}
	defer ole.CoUninitialize()

	unk, err := oleutil.CreateObject("adox.catalog")
	if err != nil {
		log.Fatal(err)
	}
	cat, err := unk.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	//OLEを使ったDB操作　（ODBCではない)
	_, err = oleutil.CallMethod(cat, "create", fmt.Sprintf("provider=Microsoft.ACE.OLEDB.12.0;data source=%s;", dbfilename))
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	createDB("./test.accdb")
	db := ConnectDB("./test.accdb")
	initDB(db)
	showTables(db)
	queryKeiyaku(db)
}
