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

func renderCodeMasterEnquiry(ctx *gin.Context) {
	var req struct {
		Domain       string `form:"Domain"       binding:"required"`
		CodeValue    string `form:"CodeValue"   binding:"required"`
		CodeName     string `form:"CodeName"   binding:"required"`
		MaxFetchRows string `form:"MaxFetchRows"   binding:"required"`
	}
	ctx.Bind(&req) //HTTP Requestのパラメータをreq構造体に紐づける
	req.Domain = rtrim(req.Domain)
	req.CodeValue = rtrim(req.CodeValue)
	req.CodeName = rtrim(req.CodeName)
	if req.MaxFetchRows == "" {
		req.MaxFetchRows = "500"
	}
	log.Print("request parameter exists:" + fmt.Sprint(req))

	list := [][]string{}
	var domains []string
	for _, v := range CodeMasterDomainList {
		if strings.Contains(v, req.Domain) {
			domains = append(domains, v)
		}
	}
	if req.Domain != "" {
		for _, v := range domains {
			for nk, nv := range CodeMaster[v] {
				if (req.CodeValue == "" || req.CodeValue == nk) &&
					(req.CodeName == "" || strings.Contains(nv, req.CodeName)) {
					list = append(list, []string{v, nk, nv})
				}
			}
		}
	} else {
		for k, v := range CodeMaster {
			for nk, nv := range v {
				if (req.CodeValue == "" || req.CodeValue == nk) &&
					(req.CodeName == "" || strings.Contains(nv, req.CodeName)) {
					list = append(list, []string{k, nk, nv})
				}
			}
		}
	}
	max, _ := strconv.Atoi(req.MaxFetchRows)
	html := ""
	if len(list) < max {
		max = len(list)
	}
	html += "<font color='grey'>" + fmt.Sprint(max) + " rows fetched</font>"
	html += "<TABLE CELLSPACING=0 CLASS='DataTable'>\n"
	html += "<TR><TH>ドメイン</TH><TH>コード値</TH><TH>コード名称</TH></TR>"

	currentRow := 0
	for _, v := range list {
		html += "<TR><TD>" + rtrim(v[0]) + "</TD><TD>" + rtrim(v[1]) + "</TD><TD>" + rtrim(v[2]) + "</TD></TR>"
		currentRow++
		if currentRow >= max {
			log.Printf("break %v %v", currentRow, max)
			break
		}
	}
	html += "</TABLE>"

	formDomain := buildInputTextField("Domain", req.Domain)
	formCodeValue := buildInputTextField("CodeValue", req.CodeValue)
	formCodeName := buildInputTextField("CodeName", req.CodeName)
	formMaxFetchRows := buildInputTextField("MaxFetchRows", req.MaxFetchRows)

	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "CodeMaster.html", gin.H{
		"formDomain":       formDomain,
		"formCodeValue":    formCodeValue,
		"formCodeName":     formCodeName,
		"formMaxFetchRows": formMaxFetchRows,
		"HTMLTable":        template.HTML(html),
	})
	log.Print("DONE: renderCodeMasterEnquiry")
}
