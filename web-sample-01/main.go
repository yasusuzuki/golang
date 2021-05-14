package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"muzzammil.xyz/jsonc"

	//SQLite
	//_ "github.com/mattn/go-sqlite3"

	//ODBCからaccdb
	//32bit の"{Microsoft Access Driver (*.mdb, *.accdb)}"しかない場合は、
	//32bitでビルドするとうまくいく。64bitアプリは64bitのODBCドライバしか使えないという制約があるため。
	//https://qiita.com/zetamatta/items/e44961a8bcbb2578cfe7
	_ "github.com/alexbrainman/odbc"

	//IBM DB2
	_ "github.com/ibmdb/go_ibm_db"
)

type htmlTableCallBack func(key string, val string, columns []string, values DBRecord) string

/**  buildHTMLTablefromDB()
* 戻り値をStringにすると、ＨＴＭＬタグもそのまま表示されてしまう。そのため、戻り値は
* template.HTML型で返す必要がある
* https://stackoverflow.com/questions/41931082/inserting-html-to-golang-template
 */
func rtrim(s interface{}) string {
	return strings.TrimSpace(fmt.Sprint(s))
}

/**
golangのstringはマルチバイト文字列に対応していないので、
rune()という機能を使う。
  https://qiita.com/hokaccha/items/3d3f45b5927b4584dbac
  https://qiita.com/reiki4040/items/b82bf5056ee747dcf713
*/
/*
func truncate(s string) string {
	r := []rune(s)
	if len(r) >= 16 {
		return string(r[0:8]) + "<br>" + string(r[8:16]) + "<br>" + string(r[16:])
	} else if len(r) >= 8 {
		return string(r[0:8]) + "<br>" + string(r[8:])
	} else {
		return s
	}
}
*/
var DBSystemColumns map[string]bool = map[string]bool{
	"証券番号枝番＿番号":        true,
	"契約管理区分キー＿英数カナ":    true,
	"契約管理レコードＩＤ＿英数カナ":  true,
	"有効開始年月日枝番＿番号":     true,
	"遡及連続＿番号":          true,
	"契約計上処理回数＿数":       true,
	"取扱年月日＿日付":         true,
	"繰越データ抽出キー＿英数カナ":   true,
	"データステータス区分＿コード":   true,
	"排他制御バージョン番号＿数":    true,
	"ビジネスタスクＩＤ＿英数カナ":   true,
	"イベント発生タイムスタンプ＿日付": true,
	"論理削除＿コード":         true,
	"データ登録タイムスタンプ＿日付":  true,
	"データ登録ユーザーＩＤ＿英数カナ": true,
	"データ登録プログラム＿英数カナ":  true,
	"データ更新タイムスタンプ＿日付":  true,
	"データ更新ユーザーＩＤ＿英数カナ": true,
	"データ更新プログラム＿英数カナ":  true,
	"リレーション用契約計上枝番＿番号": true,
}

func buildHTMLTablefromDB(sql string, callback map[string]htmlTableCallBack) template.HTML {
	columns, records, err := DBAccess(sql)
	if err != nil {
		return template.HTML(fmt.Sprintf("<pre>%s</pre><div class='message_error'>%s</div>", sql, err.Error()))
	}
	var hideDBSystemColumns bool = true
	var displaySQL bool = false
	var displayPhysicalName bool = false
	if callback["VERBOSE_MODE_FLAG"] != nil && callback["VERBOSE_MODE_FLAG"]("", "", nil, nil) == "on" {
		hideDBSystemColumns = false
		displaySQL = true
		displayPhysicalName = true
	}
	html := ""
	if displaySQL {
		html = fmt.Sprintf("<PRE>%s</PRE>", sql)
	}

	html += "<font color='grey'>" + fmt.Sprint(len(records)) + " rows fetched</font>"
	html += "<TABLE CELLSPACING=0 CLASS='DataTable'>\n"

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
		if hideDBSystemColumns && DBSystemColumns[column] {
			continue
		}
		physicalColumnName := ""
		if displayPhysicalName {
			physicalColumnName = fmt.Sprintf("<BR><FONT COLOR='lightgrey'>[%s]</FONT>", L2P(column))
		}
		if callback["H_"+column] != nil {
			html += "<TH>" + callback[column](column, "", columns, nil) + "</TH>"
		} else if strings.Contains(column, "＿コード") || strings.Contains(column, "_CD") {
			domain := findDomain(column)
			if domain != "" {
				html += "<TH STYLE='max-width:100px'><A HREF='/codeMasterEnquiry?Domain=" + domain + "'>" + column + "</A>" + physicalColumnName + "</TH>"
			} else {
				html += "<TH>" + column + physicalColumnName + "</TH>"
			}
		} else {
			if DBSystemColumns[column] {
				html += "<TH style='background-color:darkgray'>" + column + physicalColumnName + "</TH>"
			} else {
				html += "<TH>" + column + physicalColumnName + "</TH>"
			}
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
			if hideDBSystemColumns && DBSystemColumns[column] {
				continue
			}
			orig_val := record[column]
			val := rtrim(orig_val)

			if strings.Contains(column, "＿コード") {
				domain := findDomain(column)
				if domain != "" && CodeMaster[domain][val] != "" {
					val = val + "<font color='grey'>[" + CodeMaster[domain][val] + "]</font>"
				}
			}
			if callback[column] != nil {
				//DBrecordに関係なくTableに追加した列があれば表示
				html += "<TD>" + callback[column](column, val, columns, record) + "</TD>"
			} else {
				if record[column] == nil {
					html += "<TD>" + "<font color='grey'>NULL</font>" + "</TD>"
				} else {
					html += "<TD>" + val + "</TD>"
				}
			}
		}
		html += "</TR>\n"
	}
	html += "</TABLE><BR>\n"
	return template.HTML(html)
}

/*
  interface{} は原始型でも構造体でもすべてにあてはまるなんでもありの型
  https://www.tohoho-web.com/ex/golang.html#interface
  interface{}のmapはDBの１レコード分を表す。
*/
type DBRecord map[string]interface{}

func DBAccess(sqlString string) ([]string, []DBRecord, error) {
	//make([]map[string]interface{},0,5)はエラーだが、
	//なぜか以下だとうまくいく
	//    type DBRecord map[string]interface{}
	//    make([]DBRecord, 0, 5)
	// https://stackoverflow.com/questions/35362459/golang-create-a-slice-of-maps
	array := make([]DBRecord, 0, 5)
	var columns []string
	log.Print("RUN QUERY --- " + sqlString)
	rows, err := CurrentDB.DBConnection.Query(sqlString)
	if err != nil {
		return nil, nil, err
	}
	log.Print("DONE Query")
	defer rows.Close()

	columns, err = rows.Columns()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(columns); i++ {
		columns[i] = P2L(columns[i])
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
				//log.Printf("row data [%s] [%v]", col, row[i])
				row[i] = string(row[i].([]byte))
			//提案番号が0000064が64になってしまうので、以下はコメントアウト
			//num, err := strconv.Atoi(row[i].(string))
			//if err == nil {
			//	row[i] = num
			//}
			case time.Time:
				if row[i].(time.Time).Format("15:04:05") == "00:00:00" {
					row[i] = fmt.Sprint(row[i].(time.Time).Format("2006-01-02"))
				} else {
					row[i] = fmt.Sprint(row[i].(time.Time).Format("2006-01-02 15:04:05"))
				}

			}
			//log.Printf("data [%s] [%s] [%v]", col, reflect.TypeOf(row[i]), row[i])
			//TODO: 他にいい方法がないか考える
			//以下のケースで、t2.keyがNULLとなるときがある。この場合にＮＵＬＬで上書きされないように。
			//select a.*,b.* from t1 LEFT JOIN t2 where t1.key = t2.key
			if rowMap[col] == nil {
				rowMap[col] = row[i]
			}
		}
		array = append(array, rowMap)
	}

	return columns, array, nil
}

func buildInputTextField(fieldName string, value string) template.HTML {
	return template.HTML("<INPUT TYPE='TEXT' SIZE='33' CLASS='texta' NAME='" + fieldName + "' VALUE='" + value + "'></INPUT>")
}
func buildInputPullDown(fieldName string, optionValues []string, optionNames []string, selected string) template.HTML {
	html := "<SELECT NAME='" + fieldName + "' onChange=''>"
	for i, opt := range optionValues {
		if opt == selected {
			html += "<OPTION  VALUE='" + opt + "' SELECTED>" + optionNames[i] + "</OPTION>"
		} else {
			html += "<OPTION  VALUE='" + opt + "'>" + optionNames[i] + "</OPTION>"
		}
	}
	html += "</INPUT>"
	return template.HTML(html)
}
func buildInputCheckbox(fieldName string, value bool) template.HTML {
	var checkedString string
	if value {
		checkedString = "checked='checked'"
	} else {
		checkedString = ""
	}
	return template.HTML("<INPUT TYPE='CHECKBOX' NAME='" + fieldName + "' " + checkedString + " ></INPUT>")
}

var CodeMaster map[string]map[string]string
var CodeMasterDomainList []string

func findDomain(fieldName string) string {
	for _, v := range CodeMasterDomainList {
		if strings.Contains(fieldName, v) {
			return v
		}
	}
	return ""
}
func initCodeMaster() {
	CodeMaster = map[string]map[string]string{}
	f, err := os.Open(Config.CodeMasterFilePath)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//log.Print(record)
		if CodeMaster[record[0]] == nil {
			CodeMaster[record[0]] = map[string]string{record[1]: record[5]}
			CodeMasterDomainList = append(CodeMasterDomainList, record[0])
		} else {
			CodeMaster[record[0]][record[1]] = record[5]
		}
	}

}

var L2PDBTables map[string]string
var L2PDictionary map[string]string
var P2LDBTables map[string]string
var P2LDictionary map[string]string

var DBTables []string //テーブル一覧の順序を保ちたいKeiyakuEnquiryで利用する。

func L2P_SQL(sql_logical string) string {
	var multiByteCharStart bool
	var multiByteWord []rune
	var sql_physical string
	for _, v := range sql_logical {
		if v >= utf8.RuneSelf {
			//マルチバイト文字の場合
			//https://golang.hateblo.jp/entry/golang-string-byte-rune#%E3%82%B7%E3%83%B3%E3%82%B0%E3%83%AB%E3%83%90%E3%82%A4%E3%83%88%E6%96%87%E5%AD%97%E3%81%A8%E3%83%9E%E3%83%AB%E3%83%81%E3%83%90%E3%82%A4%E3%83%88%E6%96%87%E5%AD%97%E3%81%AE%E5%88%A4%E5%88%A5
			//log.Print("L2P SQL [" + string(v) + "] MULTIBYTE " + sql_physical)
			multiByteCharStart = true
			multiByteWord = append(multiByteWord, v)

		} else if multiByteCharStart && v < utf8.RuneSelf {
			//一つ前はマルチバイトで、この文字はシングルバイトの場合
			//log.Print("L2P SQL [" + string(v) + "] MULTIBYTE LAST" + sql_physical)
			multiByteCharStart = false
			sql_physical = sql_physical + L2P(string(multiByteWord)) + string(v)
			multiByteWord = []rune("")
		} else {
			//一つ前も、この文字もシングルバイトの場合
			//log.Print("L2P SQL [" + string(v) + "] SINGLEBYTE" + sql_physical)
			sql_physical += string(v)
		}
	}
	return sql_physical
}

func initL2PDictionary() {
	DBTables = make([]string, 300)
	L2PDBTables = make(map[string]string, 300)
	L2PDictionary = make(map[string]string, 20000)
	//DBテーブル一覧のロード
	f, err := os.Open(Config.DBTableListFilePath)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if L2PDBTables[record[0]] != "" {
			panic("Unexpected duplicate key in L2PDictionary[" + record[0] + "]")
		}
		//record[0]がテーブル論理名。record[1]がテーブル物理名
		L2PDBTables[record[0]] = record[1]
		DBTables = append(DBTables, record[0])
	}

	//データディクショナリのロード
	f, err = os.Open(Config.DataDictionaryFilePath)
	if err != nil {
		log.Fatal(err)
	}
	r = csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if L2PDictionary[record[0]] != "" {
			panic("Unexpected duplicate key in L2PDictionary[" + record[0] + "]")
		}
		L2PDictionary[record[0]] = record[2]
	}
	P2LDBTables = make(map[string]string, 300)
	P2LDictionary = make(map[string]string, 20000)
	for k, v := range L2PDBTables {
		P2LDBTables[v] = k
	}
	for k, v := range L2PDictionary {
		P2LDictionary[v] = k
	}

	//fmt.Printf("%#v", L2PDictionary)
}
func L2P(logicalName string) string {
	if Config.DB_SERVER_PRODUCT == "ODBC" || Config.DB_SERVER_PRODUCT == "SQLITE" {
		return logicalName
	}
	if L2PDictionary[logicalName] != "" {
		return L2PDictionary[logicalName]
	} else if L2PDBTables[logicalName] != "" {
		return L2PDBTables[logicalName]
	} else {
		return logicalName
	}
}
func P2L(physicalName string) string {
	if Config.DB_SERVER_PRODUCT == "ODBC" || Config.DB_SERVER_PRODUCT == "SQLITE" {
		return physicalName
	}
	if P2LDictionary[physicalName] != "" {
		return P2LDictionary[physicalName]
	} else if P2LDBTables[physicalName] != "" {
		return P2LDBTables[physicalName]
	} else {
		return physicalName
	}
}

var CurrentDB struct {
	Environment  string
	DBConnection *sql.DB
}

func ConnectDB(env string) error {
	var err error

	var connProperties ConnectionProperties
	connProperties = *Config.DBConnProps[0]
	for _, v := range Config.DBConnProps {
		if v.ENV == env {
			connProperties = *v
		}
	}
	log.Print("CONNECTING TO ..[" + connProperties.ENV + "]")

	var conn *sql.DB
	if Config.DB_SERVER_PRODUCT == "SQLITE" {
		conn, err = sql.Open("sqlite3", "./data/keiyaku.db")
	} else if Config.DB_SERVER_PRODUCT == "ODBC" {
		fmt.Printf("Connecting to DRIVER={Microsoft Access Driver (*.mdb, *.accdb)};DBQ=" + connProperties.DATABASE + ";\n")
		conn, err = sql.Open("odbc", "DRIVER={Microsoft Access Driver (*.mdb, *.accdb)};DBQ="+connProperties.DATABASE+";")
	} else if Config.DB_SERVER_PRODUCT == "DB2" {
		fmt.Printf("Connecting to HOSTNAME=%s;PORT=%s;DATABASE=%s;CurrentSchema=%s;UID=%s;PWD=#######\n",
			connProperties.HOSTNAME,
			connProperties.DATABASE,
			connProperties.PORT,
			connProperties.SCHEMA,
			connProperties.UID,
		)
		conn, err = sql.Open("go_ibm_db",
			fmt.Sprintf("HOSTNAME=%s;DATABASE=%s;PORT=%s;CurrentSchema=%s;UID=%s;PWD=%s",
				connProperties.HOSTNAME,
				connProperties.DATABASE,
				connProperties.PORT,
				connProperties.SCHEMA,
				connProperties.UID,
				connProperties.PWD,
			))
	}
	if err != nil {
		log.Printf("DB connection failed")
		return err
	}
	log.Print("DB CONNECTED ")

	CurrentDB.DBConnection = conn
	CurrentDB.Environment = connProperties.ENV
	return nil
}

type ConnectionProperties struct {
	ENV      string `json:"ENV"`
	HOSTNAME string `json:"HOSTNAME"`
	DATABASE string `json:"DATABASE"`
	PORT     string `json:"PORT"`
	UID      string `json:"UID"`
	PWD      string `json:"PWD"`
	SCHEMA   string `json:"SCHEMA"`
}

var Config struct {
	DB_SERVER_PRODUCT      string                  `json:"DB_SERVER_PRODUCT"`
	DBConnProps            []*ConnectionProperties `json:"DBConnection"`
	DataDictionaryFilePath string                  `json:"DataDictionaryFilePath"`
	DBTableListFilePath    string                  `json:"DBTableListFilePath"`
	CodeMasterFilePath     string                  `json:"CodeMasterFilePath"`
}

func loadConfigFile() {
	jsonStringWithComments, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	jsonString := jsonc.ToJSON(jsonStringWithComments)
	//var c Config
	if err := json.Unmarshal(jsonString, &Config); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", Config)
	fmt.Printf("%+v\n", Config.DBConnProps[0])
}

func ListEnvironment() []string {
	var dblist []string
	for _, v := range Config.DBConnProps {
		dblist = append(dblist, v.ENV)
	}
	return dblist
}

func main() {
	loadConfigFile()
	log.Print("USERDOMIN = " + os.Getenv("USERDOMAIN"))
	initL2PDictionary()
	initCodeMaster()
	dblist := ListEnvironment()
	CurrentDB.Environment = dblist[0]
	ConnectDB(CurrentDB.Environment)

	router := gin.Default()

	//gin.H{}には、原始型か構造体しか設定できないが、
	//グローバル関数はこちらで設定しておくと、テンプレート側から関数として参照できる
	//ただし、このやりかただと、すべてのテンプレートが共通の関数を使うことしかできないので
	//汎用的に設計しないといけない
	router.SetFuncMap(template.FuncMap{
		"buildHTMLTablefromDB": buildHTMLTablefromDB,
		"L2P_SQL":              L2P_SQL,
	})
	router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets")

	router.GET("/", renderKeiyakuList) //DEFAULT
	router.GET("/keiyakuList", renderKeiyakuList)
	router.GET("/keiyakuEnquiry", renderKeiyakuEnquiry)
	router.GET("/codeMasterEnquiry", renderCodeMasterEnquiry)
	router.GET("/dataDictionaryEnquiry", renderDataDictionaryEnquiry)
	router.GET("/teianList", renderTeianList)
	router.GET("/teianEnquiry", renderTeianEnquiry)

	router.Run(":8080")
	CurrentDB.DBConnection.Close()

}
