package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

var CodeMaster map[string]map[string]string
var DomainList []string

func main() {
	CodeMaster = map[string]map[string]string{}
	f, err := os.Open("./data/codemaster.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//log.Print(record)
		if CodeMaster[record[0]] == nil {
			CodeMaster[record[0]] = map[string]string{record[1]: record[5]}
			DomainList = append(DomainList, record[0])
		} else {
			CodeMaster[record[0]][record[1]] = record[5]
		}
	}

	fmt.Printf("%#v\n", DomainList)

}
