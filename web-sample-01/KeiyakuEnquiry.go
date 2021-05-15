package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func renderKeiyakuEnquiry(ctx *gin.Context) {
	//１．　HTTPリクエストパラメータを解析する
	type TableSQL struct {
		LogicalTableName string
		Sql              string
	}
	var req struct {
		PolicyNumber []string `form:"PolicyNumber"       binding:"required"`
		VerboseMode  string   `form:"VerboseMode"         binding:"required"`
	}
	ctx.Bind(&req) //HTTP Requestのパラメータをreq構造体に紐づける
	log.Printf("request parameter [%+v]", req)

	//２．　データベースに接続する
	conn, err := ConnectDB(CurrentDB.Environment)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	//defer conn.Close()　// Access DBだとここでハングアップしてしまうので、クローズしないことにした

	//３．　SQLを組み立てる
	sqls := make([]TableSQL, 0, 100)

	for _, logicalTableName := range DBTables {
		//DBテーブル一覧から契約関連エンティティだけを抽出する
		if !strings.HasPrefix(logicalTableName, "保険契約") &&
			!strings.HasPrefix(logicalTableName, "請求保険料") &&
			!strings.HasPrefix(logicalTableName, "保険対象") {
			continue
		}
		var polNums []string
		for _, v := range req.PolicyNumber {
			polNums = append(polNums, "'"+strings.Replace(rtrim(v), ",", "','", -1)+"'")
		}
		polNumsString := strings.Join(polNums, ",")
		if polNumsString == "" {
			sqls = append(sqls, TableSQL{logicalTableName, "NO_POLICY_NUMBER"})
		} else if L2P(logicalTableName) == "NO_PHYSICAL_TABLE" {
			sqls = append(sqls, TableSQL{logicalTableName, "NO_PHYSICAL_TABLE"})
		} else {
			sqls = append(sqls, TableSQL{logicalTableName, fmt.Sprintf("SELECT * FROM %s WHERE %s IN (%s) ", L2P(logicalTableName), L2P("証券＿番号"), polNumsString)})
		}
	}

	//４．　入力フォームを組み立てる
	formPolicyNumber := buildInputTextField("PolicyNumber", strings.Join(req.PolicyNumber, ","))
	formVerboseMode := buildInputCheckbox("VerboseMode", req.VerboseMode == "on")

	//５．　テンプレートエンジンからのコールバックを定義する
	callback := map[string]htmlTableCallBack{
		"VERBOSE_MODE_FLAG": func(key string, val string, columns []string, values DBRecord) string {
			return req.VerboseMode
		},
	}

	//６．　GINフレームワークのテンプレートエンジンを呼ぶ
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "KeiyakuEnquiry.html", gin.H{
		"formPolicyNumber":  formPolicyNumber,
		"htmlTableCallBack": callback,
		"conn":              conn,
		"SQLs":              sqls,
		"formVerboseMode":   formVerboseMode,
	})

	log.Print("DONE: renderKeiyakuEnquiry")

}
