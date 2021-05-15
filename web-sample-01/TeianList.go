package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderTeianList(ctx *gin.Context) {
	//１．　HTTPリクエストパラメータを解析する
	var req struct {
		AnkenNumber      string `form:"AnkenNumber"       binding:"required"`
		PolicyHolderName string `form:"PolicyHolderName"   binding:"required"`
		PolicyType       string `form:"PolicyType"         binding:"required"`
		Environment      string `form:"Environment"         binding:"required"`
		MaxFetchRows     string `form:"MaxFetchRows"         binding:"required"`
	}
	ctx.Bind(&req) //HTTP Requestのパラメータをreq構造体に紐づける

	req.AnkenNumber = rtrim(req.AnkenNumber)
	req.PolicyHolderName = rtrim(req.PolicyHolderName)
	req.PolicyType = rtrim(req.PolicyType)
	req.Environment = rtrim(req.Environment)
	if req.MaxFetchRows == "" {
		req.MaxFetchRows = "100"
	}
	log.Printf("request parameter [%+v]", req)

	//２．　データベースに接続する
	conn, err := ConnectDB(req.Environment)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	//defer conn.Close()　// Access DBだとここでハングアップしてしまうので、クローズしないことにした

	//３．　SQLを組み立てる
	sql := ""
	sql += fmt.Sprintf("SELECT a.%s,a.%s,a.%s,a.%s,d.%s,   a.%s,d.%s,d.%s,a.%s,a.%s,a.%s,   b.%s,a.%s,a.%s,a.%s,a.%s,a.%s ",
		L2P("提案案件＿番号"), L2P("提案案件番号枝番＿番号"), L2P("提案連続＿番号"), L2P("提案設計データバージョン番号＿数"), L2P("保険契約明細＿番号"),
		L2P("証券＿番号"), L2P("保険契約保険種目＿コード"), L2P("単位商品＿コード"), L2P("保険契約明細区分＿コード"), L2P("保険契約消滅変更当否＿フラグ"), L2P("契約始期年月日＿日付"),
		L2P("契約者氏名＿漢字"), L2P("団体＿コード"), L2P("代理店＿コード"), L2P("代理店サブ＿コード"), L2P("契約保険期間年＿数"), L2P("イベント発生タイムスタンプ＿日付"))
	sql += fmt.Sprintf(" FROM %s a ", L2P("提案"))
	sql += fmt.Sprintf(" INNER JOIN %s d ON a.%s = d.%s AND a.%s  = d.%s AND a.%s  = d.%s AND a.%s  = d.%s AND a.%s  = d.%s",
		L2P("提案明細"), L2P("契約管理区分キー＿英数カナ"), L2P("契約管理区分キー＿英数カナ"), L2P("提案設計データバージョン番号＿数"), L2P("提案設計データバージョン番号＿数"),
		L2P("提案案件＿番号"), L2P("提案案件＿番号"), L2P("提案案件番号枝番＿番号"), L2P("提案案件番号枝番＿番号"), L2P("提案連続＿番号"), L2P("提案連続＿番号"))
	sql += fmt.Sprintf(" INNER JOIN %s b ON a.%s = b.%s AND a.%s  = b.%s AND a.%s  = d.%s AND a.%s  = d.%s AND a.%s  = d.%s",
		L2P("提案．契約者"), L2P("契約管理区分キー＿英数カナ"), L2P("契約管理区分キー＿英数カナ"), L2P("提案設計データバージョン番号＿数"), L2P("提案設計データバージョン番号＿数"),
		L2P("提案案件＿番号"), L2P("提案案件＿番号"), L2P("提案案件番号枝番＿番号"), L2P("提案案件番号枝番＿番号"), L2P("提案連続＿番号"), L2P("提案連続＿番号"))

	if req.AnkenNumber != "" {
		sql += " AND a." + L2P("提案案件＿番号") + " LIKE '%" + req.AnkenNumber + "%' "
	}
	if req.PolicyHolderName != "" {
		sql += " AND b." + L2P("契約者氏名＿漢字") + " LIKE '%" + req.PolicyHolderName + "%' "
	}
	if req.PolicyType != "全部" && req.PolicyType != "" {
		sql += " AND d." + L2P("保険契約保険種目＿コード") + " = '" + req.PolicyType + "' "
	}

	sql += " ORDER BY a." + L2P("提案案件＿番号") + " DESC"

	sql += " FETCH FIRST " + req.MaxFetchRows + " ROWS ONLY "

	//４．　入力フォームを組み立てる
	formAnkenNumber := buildInputTextField("AnkenNumber", req.AnkenNumber)
	formPolicyHolderName := buildInputTextField("PolicyHolderName", req.PolicyHolderName)
	formPolicyType := buildInputPullDown("PolicyType", []string{"", "73", "71"}, []string{"全部", "自動車保険", "傷害保険"}, req.PolicyType)
	formMaxFetchRows := buildNumberTextField("MaxFetchRows", req.MaxFetchRows)
	formEnvironment := buildInputPullDown("Environment", listEnvironment(), listEnvironment(), CurrentDB.Environment)

	//５．　テンプレートエンジンからのコールバックを定義する
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
			html += "<INPUT TYPE='submit' NAME='ACTN' VALUE='提案DB詳細' class='button' onClick='setActionToDataForm(\"/teianEnquiry\");'>"
			return html
		},
		"提案案件＿番号": func(key string, val string, columns []string, values DBRecord) string {
			return "<a href='/teianEnquiry?AnkenNumber=" + val + "'>" + val + "</a>"
		},
		"証券＿番号": func(key string, val string, columns []string, values DBRecord) string {
			return "<a href='/keiyakuEnquiry?PolicyNumber=" + val + "'>" + val + "</a>"
		},
		"HIDE_DB_SYSTEM_COLUMNS_FLAG": func(key string, val string, columns []string, values DBRecord) string {
			return "off"
		},
	}

	//６．　GINフレームワークのテンプレートエンジンを呼ぶ
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "TeianList.html", gin.H{
		"formAnkenNumber":      formAnkenNumber,
		"formPolicyHolderName": formPolicyHolderName,
		"formPolicyType":       formPolicyType,
		"SQL":                  sql,
		"conn":                 conn,
		"htmlTableCallBack":    callback,
		"formEnvironment":      formEnvironment,
		"formMaxFetchRows":     formMaxFetchRows,
	})

	log.Print("DONE: renderTeianList")
}
