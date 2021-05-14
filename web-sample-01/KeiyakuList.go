package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderKeiyakuList(ctx *gin.Context) {
	var req struct {
		PolicyNumber     string `form:"PolicyNumber"       binding:"required"`
		PolicyHolderName string `form:"PolicyHolderName"   binding:"required"`
		PolicyType       string `form:"PolicyType"         binding:"required"`
		Environment      string `form:"Environment"         binding:"required"`
		MaxFetchRows     string `form:"MaxFetchRows"         binding:"required"`
	}
	ctx.Bind(&req) //HTTP Requestのパラメータをreq構造体に紐づける
	req.PolicyNumber = rtrim(req.PolicyNumber)
	req.PolicyHolderName = rtrim(req.PolicyHolderName)
	req.PolicyType = rtrim(req.PolicyType)
	req.Environment = rtrim(req.Environment)
	log.Printf("request parameter [%+v]", req)

	err := ConnectDB(req.Environment)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer CurrentDB.DBConnection.Close()

	sql := ""
	if Config.DB_SERVER_PRODUCT == "LOCAL_ODBC" {
		sql += "SELECT a.[証券＿番号],a.保険契約明細区分＿コード,a.保険契約消滅変更当否＿フラグ,a.[契約始期年月日＿日付],b.[契約者氏名＿漢字], c.[自動車保険契約種目ノンフリートフリート区分＿コード], c.[自動車保険契約種目フリート契約形態＿コード],a.[団体＿コード],a.[代理店＿コード],a.[代理店サブ＿コード],a.[契約保険期間年＿数]"
		sql += " FROM 保険契約 a "
		sql += " INNER JOIN ( [保険契約．契約者] b INNER JOIN [保険契約種目＿自動車] c ON b.[証券＿番号] = c.[証券＿番号] AND b.[契約計上枝番＿番号] = c.[契約計上枝番＿番号]) "
		sql += "  ON a.[証券＿番号] = b.[証券＿番号]  AND a.[契約計上枝番＿番号] = b.[契約計上枝番＿番号]"
		sql += " WHERE a.[契約計上枝番＿番号]='00001' AND b.[契約計上枝番＿番号]='00001' AND c.[契約計上枝番＿番号]='00001' AND b.[契約者ロール＿コード] = '01'   "

	} else {
		sql += fmt.Sprintf("SELECT a.%s,d.%s,a.%s,a.%s,a.%s,b.%s,a.%s,a.%s,a.%s,a.%s,a.%s", L2P("証券＿番号"), L2P("保険契約保険種目＿コード"), L2P("保険契約明細区分＿コード"), L2P("保険契約消滅変更当否＿フラグ"), L2P("契約始期年月日＿日付"),
			L2P("契約者氏名＿漢字"), L2P("団体＿コード"), L2P("代理店＿コード"), L2P("代理店サブ＿コード"), L2P("契約保険期間年＿数"), L2P("イベント発生タイムスタンプ＿日付"))
		sql += fmt.Sprintf(" FROM %s a ", L2P("保険契約"))
		sql += fmt.Sprintf(" INNER JOIN %s d ON a.%s = d.%s AND a.%s  = d.%s ", L2P("保険契約種目"), L2P("証券＿番号"), L2P("証券＿番号"), L2P("契約計上枝番＿番号"), L2P("契約計上枝番＿番号"))
		sql += fmt.Sprintf(" INNER JOIN %s b ON a.%s = b.%s AND a.%s  = b.%s ", L2P("保険契約．契約者"), L2P("証券＿番号"), L2P("証券＿番号"), L2P("契約計上枝番＿番号"), L2P("契約計上枝番＿番号"))
		//	sql += fmt.Sprintf(" LEFT JOIN %s c ON b.%s = c.%s AND b.%s = c.%s AND c.%s='00001'", L2P("保険契約種目＿自動車"), L2P("証券＿番号"), L2P("証券＿番号"), L2P("契約計上枝番＿番号"), L2P("契約計上枝番＿番号"), L2P("契約計上枝番＿番号"))
	}
	if req.PolicyNumber != "" {
		sql += " AND a." + L2P("証券＿番号") + " LIKE '%" + req.PolicyNumber + "%' "
	}
	if req.PolicyHolderName != "" {
		sql += " AND b." + L2P("契約者氏名＿漢字") + " LIKE '%" + req.PolicyHolderName + "%' "
	}
	if req.PolicyType != "全部" && req.PolicyType != "" {
		sql += " AND d." + L2P("保険契約保険種目＿コード") + " = '" + req.PolicyType + "' "
	}

	sql += " ORDER BY a." + L2P("証券＿番号") + " DESC"

	if req.MaxFetchRows == "" {
		req.MaxFetchRows = "100"
	}
	if Config.DB_SERVER_PRODUCT != "LOCAL_ODBC" {
		//FETCH FIRST doesn't work in ODBC..
		sql += " FETCH FIRST " + req.MaxFetchRows + " ROWS ONLY "
	}

	formPolicyNumber := buildInputTextField("PolicyNumber", req.PolicyNumber)
	formPolicyHolderName := buildInputTextField("PolicyHolderName", req.PolicyHolderName)
	formPolicyType := buildInputPullDown("PolicyType", []string{"", "73", "71"}, []string{"全部", "自動車保険", "傷害保険"}, req.PolicyType)
	formMaxFetchRows := buildInputTextField("MaxFetchRows", req.MaxFetchRows)
	formEnvironment := buildInputPullDown("Environment", ListEnvironment(), ListEnvironment(), CurrentDB.Environment)

	callback := map[string]htmlTableCallBack{
		"PREPEND": func(key string, val string, columns []string, values DBRecord) string {
			//注意：　key と valがブランクになるのでvaluesからpolicyNumberを取得する
			polNumColumnNumber := 0
			for i, v := range columns {
				if v == "証券＿番号" {
					polNumColumnNumber = i
					break
				}
			}
			policyNumber := fmt.Sprint(values[columns[polNumColumnNumber]])
			return "<INPUT TYPE='checkbox'  NAME='PolicyNumber' value='" + policyNumber + "'>"
		},
		"H_PREPEND": func(key string, val string, columns []string, values DBRecord) string {
			html := "<INPUT TYPE='checkbox' onClick='toggleAllMsg(this, \"PolicyNumber\");'>&nbsp;"
			html += "<INPUT TYPE='submit' NAME='ACTN' VALUE='契約DB詳細' class='button' onClick='setFormAction(\"/keiyakuEnquiry\");'>"
			return html
		},
		"証券＿番号": func(key string, val string, columns []string, values DBRecord) string {
			return "<a href='/keiyakuEnquiry?PolicyNumber=" + val + "'>" + val + "</a>"
		},
	}
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "KeiyakuList.html", gin.H{
		"formPolicyNumber":     formPolicyNumber,
		"formPolicyHolderName": formPolicyHolderName,
		"formPolicyType":       formPolicyType,
		"SQL":                  sql,
		"htmlTableCallBack":    callback,
		"formEnvironment":      formEnvironment,
		"formMaxFetchRows":     formMaxFetchRows,
	})
	log.Print("DONE: KeiyakuList.go")
}
