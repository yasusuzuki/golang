package main

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

/**
* 戻り値をStringにすると、ＨＴＭＬタグもそのまま表示されてしまう。そのため、戻り値は
* template.HTML型で返す必要がある
* https://stackoverflow.com/questions/41931082/inserting-html-to-golang-template
 */
func buildHTMLTablefromDB(sql string, value int) template.HTML {
	html := "<TABLE CELLSPACING=0 CLASS='DataTable'>\n"
	html = html + "<TR><TH></TH><TH>policyNumber</TH><TH>policyHolderName</TH><TH>policyType</TH></TR>\n"
	for _, kv := range DBAccess() {
		html = html + "<TR>"
		html = html + "<TD>" + "<INPUT TYPE='checkbox'  NAME='ref_" + fmt.Sprint(kv["policyNumber"]) + "' value='" + fmt.Sprint(kv["policyNumber"]) + "'></TD>"
		html = html + "<TD>" + fmt.Sprint(kv["policyNumber"]) + "</TD>"
		html = html + "<TD>" + fmt.Sprint(kv["policyHolderName"]) + "</TD>"
		html = html + "<TD>" + fmt.Sprint(kv["policyType"]) + "</TD>"
		html = html + "</TR>\n"
	}
	html = html + "</TABLE><BR>\n"
	return template.HTML(html)
}

func getHogeString(sql string, value int) string {
	html := "HOGE " + sql + fmt.Sprint(value)
	return html
}

/*
  interface{} は原始型でも構造体でもすべてにあてはまるなんでもありの型
  https://www.tohoho-web.com/ex/golang.html#interfaces
  interface{}のmapはDBの１レコード分を表す。
*/
type DBRecord map[string]interface{}

func DBAccess() []DBRecord {
	//make([]map[string]interface{},0,5)はエラーだが、
	//なぜか以下だとうまくいく
	//    type DBRecord map[string]interface{}
	//    make([]DBRecord, 0, 5)
	// https://stackoverflow.com/questions/35362459/golang-create-a-slice-of-maps
	array := make([]DBRecord, 0, 5)
	array = append(array, DBRecord{"policyNumber": "D00000001", "policyHolderName": "鈴木", "policyType": "ノンフリート"})
	array = append(array, DBRecord{"policyNumber": "D00000002", "policyHolderName": "田中", "policyType": "ノンフリート"})
	array = append(array, DBRecord{"policyNumber": "D00000003", "policyHolderName": "佐藤", "policyType": "フリート"})
	array = append(array, DBRecord{"policyNumber": "D00000004", "policyHolderName": "織田", "policyType": "フリート"})

	//TODO: 今は上記のようにハードコーディングしてしまっているが、今後規模が大きくなるにあたり、以下のようにファイルから抽出できるようにしたい
	/*
		https://note.crohaco.net/2019/golang-gin/
		binary, _ := ioutil.ReadFile("./users.json")
		users := make([]User, 0)
		json.Unmarshal(binary, &users)
	*/
	return array
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
	router := gin.Default()
	//gin.H{}には、原始型か構造体しか設定できないが、
	//グローバル関数はこちらで設定しておくと、テンプレート側から関数として参照できる
	//ただし、このやりかただと、すべてのテンプレートが共通の関数を使うことしかできないので管理が面倒。
	router.SetFuncMap(template.FuncMap{
		"getHogeString":        getHogeString,
		"buildHTMLTablefromDB": buildHTMLTablefromDB,
	})
	router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets")

	router.GET("/", renderKeiyakuList) //DEFAULT
	router.GET("/keiyakuList", renderKeiyakuList)
	//router.POST("/keiyakuList", renderKeiyakuList)
	router.Run(":8080")
}
