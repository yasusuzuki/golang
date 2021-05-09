package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderKeiyakuList(ctx *gin.Context) {
	var req = struct {
		policyNumber     string `form:"policyNumber"       binding:"required"`
		policyHolderName string `form:"policyHolderName"   binding:"required"`
		policyType       string `form:"policyType"         binding:"required"`
	}{ctx.Query("policyNumber"), ctx.Query("policyHolderName"), ctx.Query("policyType")}
	//}{}
	//ctx.ShouldBindWith(&req, binding.Query) //Bind() deosn't work!
	log.Print("request parameter exists:" + fmt.Sprint(req))

	sql := "SELECT * FROM 契約エンティティ WHERE 1=1 "
	if req.policyNumber != "" {
		sql += " AND policyNumber = '" + req.policyNumber + "' "
	}
	if req.policyHolderName != "" {
		sql += " AND policyHolderName = '" + req.policyHolderName + "' "
	}
	if req.policyType != "全部" && req.policyType != "" {
		sql += " AND policyType = '" + req.policyType + "' "
	}

	formPolicyNumber := buildInputTextField("policyNumber", req.policyNumber)
	formPolicyHolderName := buildInputTextField("policyHolderName", req.policyHolderName)
	formPolicyType := buildInputPullDown("policyType", []string{"全部", "ノンフリート", "フリート"}, req.policyType)
	callback := map[string]htmlTableCallBack{
		"PREPEND": func(key string, val string, columns []string, values DBRecord) string {
			//注意：　key と valがブランクになるのでvaluesからpolicyNumberを取得する
			policyNumber := fmt.Sprint(values[columns[0]])
			//return "<INPUT TYPE='checkbox'  NAME='polNum_" + policyNumber + "' value='" + policyNumber + "'>"
			return "<INPUT TYPE='checkbox'  NAME='policyNumber' value='" + policyNumber + "'>"
		},
		"H_PREPEND": func(key string, val string, columns []string, values DBRecord) string {
			html := "<INPUT TYPE='checkbox' onClick='toggleAllMsg(this, \"policyNumber\");'>&nbsp;"
			html += "<INPUT TYPE='submit' NAME='ACTN' VALUE='契約DB詳細' class='button' onClick='setFormAction(\"/keiyakuEnquiry\");'>"
			return html
		},
		"policyNumber": func(key string, val string, columns []string, values DBRecord) string {
			return "<a href='/keiyakuEnquiry?policyNumber=" + val + "'>" + val + "</a>"
		},
		"policyHolderName": func(key string, val string, columns []string, values DBRecord) string {
			return "<font color='red'>" + val + "</font>"
		},
	}
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "KeiyakuList.html", gin.H{
		"formPolicyNumber":     formPolicyNumber,
		"formPolicyHolderName": formPolicyHolderName,
		"formPolicyType":       formPolicyType,
		"SQL":                  sql,
		"htmlTableCallBack":    callback,
	})
}
