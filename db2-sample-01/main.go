package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/ibmdb/go_ibm_db"
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
	sql := "select * from YASU.TANNPO"
	//sql := "select * from POL.TB000570"
	_, records := queryForRows(db, sql)
	fmt.Printf("%#v\n", records)
}

func showTables(db *sql.DB) {
	sql := "SELECT name FROM sqlite_master" // WHERE type='table'"
	_, records := queryForRows(db, sql)
	fmt.Printf("%#v\n", records)
}
func initDB(db *sql.DB) {

}

func ConnectDB(dbfilename string) *sql.DB {
	//config := "HOSTNAME=10.240.30.11;DATABASE=LOCDBC6;PORT=456;UID=CS11146;PWD=Welcome4"
	config := "HOSTNAME=localhost;DATABASE=YASU_DB;UID=db2admin;PWD=db2admin"
	db, err := sql.Open("go_ibm_db", config)
	if err != nil {
		fmt.Printf("DBとの接続に失敗しました。%+v", err)
	}
	return db
}

func createDB(dbfilename string) {
	//TODO
}

func main() {
	db := ConnectDB("YASU_DB")
	queryKeiyaku(db)
	fmt.Printf("DONE")
}
