package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func renderDataDictionaryEnquiry(ctx *gin.Context) {
	var req struct {
		Field                         string `form:"Field"       binding:"required"`
		FieldPhysicalNameForDataModel string `form:"FieldPhysicalNameForDataModel"   binding:"required"`
		MaxFetchRows                  string `form:"MaxFetchRows"   binding:"required"`
	}
	ctx.Bind(&req) //HTTP Requestのパラメータをreq構造体に紐づける
	req.Field = rtrim(req.Field)
	req.FieldPhysicalNameForDataModel = rtrim(req.FieldPhysicalNameForDataModel)
	log.Print("request parameter exists:" + fmt.Sprint(req))

	list := [][]string{}
	if req.MaxFetchRows == "" {
		req.MaxFetchRows = "500"
	}
	max, _ := strconv.Atoi(req.MaxFetchRows)
	currentRow := 0
	for k, v := range L2PDictionary {
		if (req.Field == "" || strings.Contains(k, req.Field)) &&
			(req.FieldPhysicalNameForDataModel == "" || strings.Contains(v, req.FieldPhysicalNameForDataModel)) {
			list = append(list, []string{k, v})
			currentRow++
		}
		if currentRow >= max {
			break
		}
	}
	html := ""
	html += "<font color='grey'>" + fmt.Sprint(len(list)) + " rows fetched</font>"
	html += "<TABLE CELLSPACING=0 CLASS='DataTable'>\n"
	html += "<TR><TH>データ項目名</TH><TH>項目物理名（データモデリング）</TH></TR>"
	for _, v := range list {
		html += "<TR><TD>" + rtrim(v[0]) + "</TD><TD>" + rtrim(v[1]) + "</TD></TR>"
	}
	html += "</TABLE>"

	formField := buildInputTextField("Field", req.Field)
	formFieldPhysicalNameForDataModel := buildInputTextField("FieldPhysicalNameForDataModel", req.FieldPhysicalNameForDataModel)
	formMaxFetchRows := buildInputTextField("MaxFetchRows", req.MaxFetchRows)

	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "DataDictionaryEnquiry.html", gin.H{
		"formField":                         formField,
		"formFieldPhysicalNameForDataModel": formFieldPhysicalNameForDataModel,
		"formMaxFetchRows":                  formMaxFetchRows,
		"HTMLTable":                         template.HTML(html),
	})
	log.Print("DONE: renderDataDictionaryEnquiry")
}
