//Go source to generate custom bundled-fonts.go
//Replace...)\fyne-io\fyne\theme\bundled-fonts.go by the generated file
package main

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    "time"
    "unsafe"
)

// func ttf_to_data(filename string) string {
//  x, errx := ioutil.ReadFile(filename)
//  if errx != nil {
//      panic("Read Error in " + filename)
//  }
//  fmt.Println("Reading " + filename)

//  str := "[]byte{"
//  for _, byte1 := range x {
//      str = str + strconv.Itoa(int(byte1)) + ","
//  }
//  return strings.TrimRight(str, ",") + "}"
// }

func ttf_to_data(filename string) string {
    x, errx := ioutil.ReadFile(filename)
    if errx != nil {
        panic("Read Error in " + filename)
    }
    fmt.Println("Reading " + filename + ": length was " + strconv.Itoa(len(x)))

    vec := make([]string, len(x))
    for ix, byte1 := range x {
        vec[ix] = strconv.Itoa(int(byte1))
    }
    //fmt.Println(strings.Join(vec, ","))
    return "[]byte{" + strings.Join(vec, ",") + "}"
}

var (
    filenameRegular    = "meiryo001.ttf"    //Regular Font TTF; Place it in the same folder
    filenameBold       = "meiryo001.ttf"       //Bold
    filenameItalic     = "meiryo001.ttf"     //Italic
    filenameBoldItalic = "meiryo001.ttf" //BoldItalic
    filenameMono       = "meiryo001.ttf"    //Mono
)

func main() {
    t0 := time.Now()
    str := strings.Join([]string{
        "// **** THIS FILE IS AUTO-GENERATED Custom bundled-fonts.go//\n//Replace (gopath)\\src\\github.com\\fyne-io\\fyne\\theme\\bundled-fonts.go by the generated file//\n\npackage theme\n\nimport \"fyne.io/fyne\"\n\nvar regular = &fyne.StaticResource{\n    StaticName: \"Custom-Regular.ttf\",\n   StaticContent: ",
        ttf_to_data(filenameRegular),
        "}\nvar bold = &fyne.StaticResource{\n  StaticName: \"Custom-Bold.ttf\",\n  StaticContent: ",
        ttf_to_data(filenameBold),
        "}\nvar italic = &fyne.StaticResource{\n    StaticName: \"Custom-Italic.ttf\",\n    StaticContent: ",
        ttf_to_data(filenameItalic),
        "}\nvar bolditalic = &fyne.StaticResource{\n    StaticName: \"Custom-BoldItalic.ttf\",\n    StaticContent: ",
        ttf_to_data(filenameBoldItalic),
        "}\nvar monospace = &fyne.StaticResource{\n StaticName: \"Custom-Mono-Regular.ttf\",\n  StaticContent: ",
        ttf_to_data(filenameMono),
        "}",
    }, "")

    data := *(*[]byte)(unsafe.Pointer(&str))
    errx := ioutil.WriteFile("bundled-fonts.go", data, 0777)
    if errx != nil {
        panic("Write error")
    }

    fmt.Print("Total time was ")
    fmt.Print(time.Since(t0).Minutes())
    fmt.Print(" minutes. ")

}