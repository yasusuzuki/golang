package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello") // ウィンドウタイトル
	//myApp.Settings().SetTheme(&myTheme{})
	myWindow.SetContent(widget.NewLabel("Hello World!あいうえお")) // 文字入りラベルをウィンドウコンテンツに配置

	myWindow.ShowAndRun() // アプリケーションの実行(Show & Run)
	// ここ以降に処理を記述する場合、アプリケーションの実行が終わるまで実行されない。
}
