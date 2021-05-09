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

	sql := "SELECT * FROM KEIYAKUDB WHERE 1=1 "
	if req.policyNumber != "" {
		sql += " AND POLICY_NUMBER = '" + req.policyNumber + "' "
	}
	if req.policyHolderName != "" {
		sql += " AND POLICY_HOLDER_NAME = '" + req.policyHolderName + "' "
	}
	if req.policyType != "全部" && req.policyType != "" {
		sql += " AND POLICY_TYPE = '" + req.policyType + "' "
	}

	formPolicyNumber := buildInputTextField("policyNumber", req.policyNumber)
	formPolicyHolderName := buildInputTextField("policyHolderName", req.policyHolderName)
	formPolicyType := buildInputPullDown("policyType", []string{"全部", "ノンフリート", "フリート"}, req.policyType)
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "KeiyakuList.html", gin.H{
		"formPolicyNumber":     formPolicyNumber,
		"formPolicyHolderName": formPolicyHolderName,
		"formPolicyType":       formPolicyType,
		"SQL":                  sql,
	})
}
