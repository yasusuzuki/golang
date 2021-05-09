package main

import (
	"database/sql"
	"fmt"
	"os"
)

var KEIYAKU_DATA []DBRecord
var MEISAI_DATA []DBRecord
var TANPO_DATA []DBRecord

func initDB() {
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"policyNumber": "D00000001", "policyHolderName": "鈴木", "policyType": "ノンフリート"})
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"policyNumber": "D00000002", "policyHolderName": "田中", "policyType": "ノンフリート"})
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"policyNumber": "D00000003", "policyHolderName": "佐藤", "policyType": "フリート"})
	KEIYAKU_DATA = append(KEIYAKU_DATA, DBRecord{"policyNumber": "D00000004", "policyHolderName": "織田", "policyType": "フリート"})
	MEISAI_DATA = append(MEISAI_DATA, DBRecord{"policyNumber": "D00000002", "meisaiNumber": "0001", "carType": "ビッツ"})
	MEISAI_DATA = append(MEISAI_DATA, DBRecord{"policyNumber": "D00000003", "meisaiNumber": "0001", "carType": "BMW"})
	MEISAI_DATA = append(MEISAI_DATA, DBRecord{"policyNumber": "D00000003", "meisaiNumber": "0002", "carType": "Audi"})
	MEISAI_DATA = append(MEISAI_DATA, DBRecord{"policyNumber": "D00000004", "meisaiNumber": "0001", "carType": "AA"})
	MEISAI_DATA = append(MEISAI_DATA, DBRecord{"policyNumber": "D00000004", "meisaiNumber": "0002", "carType": "BB"})

	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000001", "meisaiNumber": "0001", "tokuyakuCode": "T001", "hokenKingaku": "0010000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000001", "meisaiNumber": "0001", "tokuyakuCode": "T002", "hokenKingaku": "0010000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000001", "meisaiNumber": "0001", "tokuyakuCode": "T003", "hokenKingaku": "0010000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000002", "meisaiNumber": "0001", "tokuyakuCode": "T001", "hokenKingaku": "0020000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000002", "meisaiNumber": "0001", "tokuyakuCode": "T002", "hokenKingaku": "0020000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000002", "meisaiNumber": "0001", "tokuyakuCode": "T003", "hokenKingaku": "0020000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000003", "meisaiNumber": "0001", "tokuyakuCode": "T001", "hokenKingaku": "0030000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000003", "meisaiNumber": "0001", "tokuyakuCode": "T002", "hokenKingaku": "0030000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000003", "meisaiNumber": "0002", "tokuyakuCode": "T001", "hokenKingaku": "0040000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000003", "meisaiNumber": "0002", "tokuyakuCode": "T002", "hokenKingaku": "0040000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000004", "meisaiNumber": "0001", "tokuyakuCode": "T001", "hokenKingaku": "0050000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000004", "meisaiNumber": "0001", "tokuyakuCode": "T002", "hokenKingaku": "0050000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000004", "meisaiNumber": "0002", "tokuyakuCode": "T001", "hokenKingaku": "0060000000"})
	TANPO_DATA = append(TANPO_DATA, DBRecord{"policyNumber": "D00000004", "meisaiNumber": "0002", "tokuyakuCode": "T002", "hokenKingaku": "0060000000"})

	var dbfile string = "./test.db"
	os.Remove(dbfile)
	//	db, err := sql.Open("sqlite3", ":memory:")
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE "契約エンティティ" ("policyNumber" CHAR(12) PRIMARY KEY, "policyHolderName" VARCHAR(255), "policyType" VARCHAR(255))`)
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare(`INSERT INTO "契約エンティティ" ("policyNumber", "policyHolderName","policyType") VALUES (?, ?,?) `)
	if err != nil {
		panic(err)
	}
	for _, v := range KEIYAKU_DATA {
		if _, err = stmt.Exec(v["policyNumber"], v["policyHolderName"], v["policyType"]); err != nil {
			panic(err)
		}
	}

	_, err = db.Exec(`CREATE TABLE "明細エンティティ" ("policyNumber" CHAR(12) , "meisaiNumber" CHAR(7) , "carType" VARCHAR(255))`)
	if err != nil {
		panic(err)
	}
	stmt, err = db.Prepare(`INSERT INTO "明細エンティティ" ("policyNumber", "meisaiNumber","carType") VALUES (?, ?,?) `)
	if err != nil {
		panic(err)
	}

	for _, v := range MEISAI_DATA {
		if _, err = stmt.Exec(v["policyNumber"], v["meisaiNumber"], v["carType"]); err != nil {
			panic(err)
		}
	}

	_, err = db.Exec(`CREATE TABLE "担保エンティティ" ("policyNumber" CHAR(12) , "meisaiNumber" CHAR(7) , "tokuyakuCode" VARCHAR(255) , "hokenKingaku" VARCHAR(255))`)
	if err != nil {
		panic(err)
	}
	stmt, err = db.Prepare(`INSERT INTO "担保エンティティ" ("policyNumber", "meisaiNumber","tokuyakuCode","hokenKingaku") VALUES (?, ?,?,?) `)
	if err != nil {
		panic(err)
	}
	for _, v := range TANPO_DATA {
		if _, err = stmt.Exec(v["policyNumber"], v["meisaiNumber"], v["tokuyakuCode"], v["hokenKingaku"]); err != nil {
			panic(err)
		}
	}
	stmt.Close()
	db.Close()
}

func testSQL() {

	var policyNumber string
	var policyHolderName string
	var policyType string

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query(`SELECT * FROM "契約エンティティ"`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&policyNumber, &policyHolderName, &policyType)
		if err != nil {
			panic(err)
		}
		fmt.Println(policyNumber, " ", policyHolderName, " ", policyType)
	}

	db.Close()
}
