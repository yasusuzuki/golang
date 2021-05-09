package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func renderKeiyakuEnquiry(ctx *gin.Context) {
	var req = struct {
		policyNumber []string `form:"policyNumber"       binding:"required"`
	}{}
	req.policyNumber, _ = ctx.GetQueryArray("policyNumber")
	//}{}
	//ctx.ShouldBindWith(&req, binding.Query) //Bind() deosn't work!
	log.Print("request parameter exists:" + fmt.Sprint(req))

	sqlKeiyaku := "SELECT * FROM 契約エンティティ WHERE 1=1 "
	if len(req.policyNumber) > 0 {
		sqlKeiyaku += " AND policyNumber IN ('" + strings.Join(req.policyNumber, "','") + "') "
	}
	sqlMeisai := "SELECT * FROM 明細エンティティ WHERE 1=1 "
	if len(req.policyNumber) > 0 {
		sqlMeisai += " AND policyNumber IN ('" + strings.Join(req.policyNumber, "','") + "') "
	}
	sqlTanpo := "SELECT * FROM 担保エンティティ WHERE 1=1 "
	if len(req.policyNumber) > 0 {
		sqlTanpo += " AND policyNumber IN ('" + strings.Join(req.policyNumber, "','") + "') "
	}
	formPolicyNumber := buildInputTextField("policyNumber", strings.Join(req.policyNumber, ","))

	callback := map[string]htmlTableCallBack{
		"policyNumber": func(key string, val string, columns []string, values DBRecord) string {
			return "<a href='/keiyakuEnquiry?policyNumber=" + val + "'>" + val + "</a>"
		},
		"policyHolderName": func(key string, val string, columns []string, values DBRecord) string {
			return "<font color='red'>" + val + "</font>"
		},
	}
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "KeiyakuEnquiry.html", gin.H{
		"formPolicyNumber":  formPolicyNumber,
		"SQLKeiyaku":        sqlKeiyaku,
		"SQLMeisai":         sqlMeisai,
		"SQLTanpo":          sqlTanpo,
		"htmlTableCallBack": callback,
	})
}
