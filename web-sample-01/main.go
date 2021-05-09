package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

type htmlTableCallBack func(key string, val string, columns []string, values DBRecord) string

/**  buildHTMLTablefromDB()
* 戻り値をStringにすると、ＨＴＭＬタグもそのまま表示されてしまう。そのため、戻り値は
* template.HTML型で返す必要がある
* https://stackoverflow.com/questions/41931082/inserting-html-to-golang-template
 */
func buildHTMLTablefromDB(sql string, callback map[string]htmlTableCallBack) template.HTML {
	html := "<TABLE CELLSPACING=0 CLASS='DataTable'>\n"
	columns, records := DBAccess(sql)

	html += "<TR>"
	//DBrecordに存在していない項目をTableに追加したい場合
	if callback["PREPEND"] != nil {
		if callback["H_PREPEND"] != nil {
			html += "<TH>" + callback["H_PREPEND"]("", "", columns, nil) + "</TH>"
		} else {
			html += "<TH></TH>"
		}
	}
	for _, column := range columns {
		if callback["H_"+column] != nil {
			//DBrecordに関係なくTableに追加した列があれば表示
			html += "<TH>" + callback[column](column, "", columns, nil) + "</TH>"
		} else {
			html += "<TH>" + fmt.Sprint(column) + "</TH>"
		}
	}
	html += "</TR>\n"

	for _, record := range records {
		html += "<TR>"
		//DBrecordに存在していない項目をTableに追加したい場合
		if callback["PREPEND"] != nil {
			html += "<TD>" + callback["PREPEND"]("", "", columns, record) + "</TD>"
		}
		for _, column := range columns {
			if callback[column] != nil {
				//DBrecordに関係なくTableに追加した列があれば表示
				html += "<TD>" + callback[column](column, fmt.Sprint(record[column]), columns, record) + "</TD>"
			} else {
				html += "<TD>" + fmt.Sprint(record[column]) + "</TD>"
			}

		}
		html += "</TR>\n"
	}
	html += "</TABLE><BR>\n"
	return template.HTML(html)
}

/*
  interface{} は原始型でも構造体でもすべてにあてはまるなんでもありの型
  https://www.tohoho-web.com/ex/golang.html#interfaces
  interface{}のmapはDBの１レコード分を表す。
*/
type DBRecord map[string]interface{}

func DBAccess(sqlString string) ([]string, []DBRecord) {
	//make([]map[string]interface{},0,5)はエラーだが、
	//なぜか以下だとうまくいく
	//    type DBRecord map[string]interface{}
	//    make([]DBRecord, 0, 5)
	// https://stackoverflow.com/questions/35362459/golang-create-a-slice-of-maps
	array := make([]DBRecord, 0, 5)
	var columns []string

	if strings.Contains(sqlString, "FROM ") {
		db, err := sql.Open("sqlite3", "./test.db")
		if err != nil {
			panic(err)
		}

		rows, err := db.Query(sqlString)
		if err != nil {
			panic(err)
		}
		columns, err = rows.Columns()
		if err != nil {
			panic(err)
		}
		for rows.Next() {
			var row = make([]interface{}, len(columns))
			var rowp = make([]interface{}, len(columns))
			for i := 0; i < len(columns); i++ {
				rowp[i] = &row[i]
			}

			rows.Scan(rowp...)

			rowMap := make(map[string]interface{})
			for i, col := range columns {
				switch row[i].(type) {
				case []byte:
					row[i] = string(row[i].([]byte))
					num, err := strconv.Atoi(row[i].(string))
					if err == nil {
						row[i] = num
					}
				}
				rowMap[col] = row[i]
			}
			array = append(array, rowMap)
		}

		db.Close()
		//array = KEIYAKU_DATA
		//columns = []string{"policyNumber", "policyHolderName", "policyType"}
	} else if strings.Contains(sqlString, "FROM 明細エンティティ") {
		array = MEISAI_DATA
		columns = []string{"policyNumber", "meisaiNumber", "carType"}
	} else if strings.Contains(sqlString, "FROM 担保特約エンティティ") {
		array = TANPO_DATA
		columns = []string{"policyNumber", "meisaiNumber", "tokuyakuCode", "hokenKingaku"}
	}
	//TODO: 今は上記のようにハードコーディングしてしまっているが、今後規模が大きくなるにあたり、以下のようにファイルから抽出できるようにしたい
	/*
		https://note.crohaco.net/2019/golang-gin/
		binary, _ := ioutil.ReadFile("./users.json")
		users := make([]User, 0)
		json.Unmarshal(binary, &users)
	*/
	return columns, array
}

func buildInputTextField(fieldName string, value string) template.HTML {
	return template.HTML("<INPUT TYPE='TEXT' SIZE='33' CLASS='texta' NAME='" + fieldName + "' VALUE='" + value + "'></INPUT>")
}
func buildInputPullDown(fieldName string, options []string, selected string) template.HTML {
	html := "<SELECT NAME='" + fieldName + "' onChange=''>"
	for _, opt := range options {
		if opt == selected {
			html += "<OPTION  VALUE='" + opt + "' SELECTED>" + opt + "</OPTION>"
		} else {
			html += "<OPTION  VALUE='" + opt + "'>" + opt + "</OPTION>"
		}
	}
	html += "</INPUT>"
	return template.HTML(html)
}

func main() {
	//initDB()
	//runSQLite()

	router := gin.Default()

	//gin.H{}には、原始型か構造体しか設定できないが、
	//グローバル関数はこちらで設定しておくと、テンプレート側から関数として参照できる
	//ただし、このやりかただと、すべてのテンプレートが共通の関数を使うことしかできないので
	//汎用的に設計しないといけない
	router.SetFuncMap(template.FuncMap{
		"buildHTMLTablefromDB": buildHTMLTablefromDB,
	})
	router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets")

	router.GET("/", renderKeiyakuList) //DEFAULT
	router.GET("/keiyakuList", renderKeiyakuList)
	router.GET("/keiyakuEnquiry", renderKeiyakuEnquiry)
	router.Run(":8080")
}
