package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func renderTeianEnquiry(ctx *gin.Context) {
	type TableSQL struct {
		LogicalTableName string
		Sql              string
	}
	var req struct {
		AnkenNumber []string `form:"AnkenNumber"       binding:"required"`
		VerboseMode string   `form:"VerboseMode"         binding:"required"`
	}
	ctx.Bind(&req) //HTTP Requestのパラメータをreq構造体に紐づける
	log.Printf("request parameter [%+v]", req)

	err := ConnectDB(CurrentDB.Environment)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	defer CurrentDB.DBConnection.Close()

	sqls := make([]TableSQL, 0, 100)

	for _, logicalTableName := range DBTables {
		//DBテーブル一覧から提案関連エンティティだけを抽出する
		if !strings.HasPrefix(logicalTableName, "提案") {
			continue
		}
		var ankenNums []string
		for _, v := range req.AnkenNumber {
			ankenNums = append(ankenNums, "'"+strings.Replace(rtrim(v), ",", "','", -1)+"'")
		}
		ankenNumsString := strings.Join(ankenNums, ",")
		if ankenNumsString == "" {
			sqls = append(sqls, TableSQL{logicalTableName, "NO_ANKEN_NUMBER"})
		} else if L2P(logicalTableName) == "NO_PHYSICAL_TABLE" {
			sqls = append(sqls, TableSQL{logicalTableName, "NO_PHYSICAL_TABLE"})
		} else {
			sqls = append(sqls, TableSQL{logicalTableName, fmt.Sprintf("SELECT * FROM %s WHERE %s IN (%s) ", L2P(logicalTableName), L2P("提案案件＿番号"), ankenNumsString)})
		}
	}
	formAnkenNumber := buildInputTextField("AnkenNumber", strings.Join(req.AnkenNumber, ","))
	formVerboseMode := buildInputCheckbox("VerboseMode", req.VerboseMode == "on")

	callback := map[string]htmlTableCallBack{
		"VERBOSE_MODE_FLAG": func(key string, val string, columns []string, values DBRecord) string {
			return req.VerboseMode
		},
	}
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "TeianEnquiry.html", gin.H{
		"formAnkenNumber":   formAnkenNumber,
		"htmlTableCallBack": callback,
		"SQLs":              sqls,
		"formVerboseMode":   formVerboseMode,
	})
	log.Print("DONE: renderKeiyakuEnquiry")

}
