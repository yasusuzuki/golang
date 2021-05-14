package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderTeianList(ctx *gin.Context) {
	var req struct {
		AnkenNumber      string `form:"AnkenNumber"       binding:"-"`
		PolicyHolderName string `form:"PolicyHolderName"   binding:"-"`
		PolicyType       string `form:"PolicyType"         binding:"-"`
		Environment      string `form:"Environment"         binding:"-"`
		MaxFetchRows     string `form:"MaxFetchRows"         binding:"-"`
	}
	ctx.Bind(&req) //HTTP Requestのパラメータをreq構造体に紐づける

	req.AnkenNumber = rtrim(req.AnkenNumber)
	req.PolicyHolderName = rtrim(req.PolicyHolderName)
	req.PolicyType = rtrim(req.PolicyType)
	req.Environment = rtrim(req.Environment)
	log.Printf("request parameter [%+v]", req)

	err := ConnectDB(req.Environment)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer CurrentDB.DBConnection.Close()

	if req.MaxFetchRows == "" {
		req.MaxFetchRows = "100"
	}

	formAnkenNumber := buildInputTextField("AnkenNumber", req.AnkenNumber)
	formPolicyHolderName := buildInputTextField("PolicyHolderName", req.PolicyHolderName)
	formPolicyType := buildInputPullDown("PolicyType", []string{"", "73", "71"}, []string{"全部", "自動車保険", "傷害保険"}, req.PolicyType)
	formMaxFetchRows := buildInputTextField("MaxFetchRows", req.MaxFetchRows)
	formEnvironment := buildInputPullDown("Environment", ListEnvironment(), ListEnvironment(), CurrentDB.Environment)

	callback := map[string]htmlTableCallBack{
		"PREPEND": func(key string, val string, columns []string, values DBRecord) string {
			//注意：　key と valがブランクになるのでvaluesからpolicyNumberを取得する
			ankenNumColumnNumber := 0
			for i, v := range columns {
				if v == "提案案件＿番号" {
					ankenNumColumnNumber = i
					break
				}
			}
			ankenNumber := fmt.Sprint(values[columns[ankenNumColumnNumber]])
			return "<INPUT TYPE='checkbox'  NAME='AnkenNumber' value='" + ankenNumber + "'>"
		},
		"H_PREPEND": func(key string, val string, columns []string, values DBRecord) string {
			html := "<INPUT TYPE='checkbox' onClick='toggleAllMsg(this, \"AnkenNumber\");'>&nbsp;"
			html += "<INPUT TYPE='submit' NAME='ACTN' VALUE='提案DB詳細' class='button' onClick='setFormAction(\"/teianEnquiry\");'>"
			return html
		},
		"提案案件＿番号": func(key string, val string, columns []string, values DBRecord) string {
			return "<a href='/teianEnquiry?AnkenNumber=" + val + "'>" + val + "</a>"
		},
		"証券＿番号": func(key string, val string, columns []string, values DBRecord) string {
			return "<a href='/keiyakuEnquiry?PolicyNumber=" + val + "'>" + val + "</a>"
		},
	}
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "TeianList.html", gin.H{
		"formAnkenNumber":      formAnkenNumber,
		"formPolicyHolderName": formPolicyHolderName,
		"formPolicyType":       formPolicyType,
		"htmlTableCallBack":    callback,
		"formEnvironment":      formEnvironment,
		"formMaxFetchRows":     formMaxFetchRows,
		"req":                  req,
	})

	log.Print("DONE: renderTeianList")
}
