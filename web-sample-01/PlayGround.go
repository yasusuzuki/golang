package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func renderPlayGround(ctx *gin.Context) {
	req := make(map[string]string)
	req["Environment"] = ctx.Query("Environment")
	var formBuilder FormBuilder
	formBuilder.ctx = ctx
	formBuilder.req = &req

	log.Printf("request parameter [%+v]", req)

	conn, err := ConnectDB(req["Environment"])
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	callback := map[string]htmlTableCallBack{
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
	//gin.H内に、キーバリュー形式の値を設定しておくと、テンプレート側から変数として参照できる　{{.変数名}}といった感じ
	ctx.HTML(http.StatusOK, "PlayGround.html", gin.H{
		"htmlTableCallBack": callback,
		"req":               formBuilder.req,
		"formBuiler":        &formBuilder, //ここはポインタを渡す
		"conn":              conn,
		//そうでないと、このエラーがでる：　template: PlayGround.html:25:52: executing "PlayGround.html" at <.formBuiler.Build>: can't evaluate field Build in type interface {}
	})

	//log.Print("Connection Closing..")
	//conn.Close()  // Access DBだとここでハングアップしてしまうので、クローズしないことにした
	//log.Print("Connection Closed..")
	log.Print("DONE: renderPlayGround")
}

type FormBuilder struct {
	ctx *gin.Context
	req *map[string]string
}

//実験場用のフォーム作成関数
func (p *FormBuilder) Build(formType string, fieldName string, options ...string) template.HTML {
	/*GINフレームワークが異常終了しないためにreqマップのキーに対して、何かしらの値を設定して初期化する必要がある。ブランクで初期化してもよい*/
	requestParamValue := p.ctx.Query(fieldName)
	//表のセルからマウスを使ってコピペすると、どうしてもタブが入ってしまう。マウス操作を簡単にするため、両端のスペースやタブは取り除く。
	(*p.req)[fieldName] = trim(requestParamValue)
	if formType == "textForm" {
		if len(options) > 1 {
			return template.HTML(fmt.Sprintf("<SPAN CLASS='error_message'>テキスト入力フォームのoptionsパラメータ数は０か１でなければいけません 設定値[%+v]</DIV>", options))
		}
		//options[0]はデフォルト値。画面入力がなければ、デフォルト値を使う
		if len(options) == 1 && requestParamValue == "" {
			(*p.req)[fieldName] = options[0]
			return buildInputTextField(fieldName, options[0])
		} else {
			return buildInputTextField(fieldName, requestParamValue)
		}
	} else if formType == "numberForm" {
		if len(options) > 1 {
			return template.HTML(fmt.Sprintf("<SPAN CLASS='error_message'>テキスト入力フォームのoptionsパラメータ数は０か１でなければいけません 設定値[%+v]</DIV>", options))
		}
		//options[0]はデフォルト値。画面入力がなければ、デフォルト値を使う
		if len(options) == 1 && requestParamValue == "" {
			(*p.req)[fieldName] = options[0]
			return buildNumberTextField(fieldName, options[0])
		} else {
			return buildNumberTextField(fieldName, requestParamValue)
		}
	} else if formType == "optionForm" {
		if len(options) != 2 {
			return template.HTML(fmt.Sprintf("<SPAN CLASS='error_message'>選択入力フォームのoptionsパラメータは２つでなければいけません 設定値[%+v]</DIV>", options))
		}
		optionValuesArray := strings.Split(options[0], ",")
		for i, v := range optionValuesArray {
			optionValuesArray[i] = trim(v)
		}
		log.Printf("optionValues %v", optionValuesArray)
		optionNamesArray := strings.Split(options[1], ",")
		for i, v := range optionNamesArray {
			optionNamesArray[i] = trim(v)
		}
		return buildInputPullDown(fieldName, optionValuesArray, optionNamesArray, requestParamValue)
	} else if formType == "checkbox" {
		return template.HTML("NULL")
	} else if formType == "environment" {
		if len(options) != 0 {
			return template.HTML(fmt.Sprintf("<SPAN CLASS='error_message'>テキスト入力フォームのoptionsパラメータは０でなければいけません 設定値[%+v]</DIV>", options))
		}
		return buildInputPullDown("Environment", listEnvironment(), listEnvironment(), CurrentDB.Environment)
	} else {
		return template.HTML("NULL")
	}
}
