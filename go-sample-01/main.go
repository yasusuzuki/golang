package main

import (
	"fmt"
	"strings"
)

func rtrim(s interface{}) string {
	return strings.TrimSpace(fmt.Sprint(s))
}

func main() {
	//var a string
	//fmt.Printf("string == nil [%v]\n", a == nil)  //invalid operation: a == nil
	var b interface{}
	fmt.Printf("interface{} is nil ? [%v]\n", b == nil)
	fmt.Printf("nil interface{} is comparable with 'aa' ? [%v]\n", b == "aa")

	var c []string
	fmt.Printf("[]string == nil [%v]\n", c == nil)
	var d []string = nil
	fmt.Printf("[]string == nil [%v]\n", d == nil)
	var e []string = []string{}
	fmt.Printf("var e string [%v]\n", e == nil)

	type ST struct {
		f1 string
	}
	//var s1 ST
	//fmt.Printf("struct == nil? [%v]\n", s1 == nil) //invalid operation: s1 == nil (mismatched types struct { f1 string } and nil)

	var s2 *ST
	fmt.Printf("pointer is null? [%v]\n", s2 == nil)
	//fmt.Printf("pointer to struct is null? [%v]\n", *s2 == nil) //invalid operation: *s2 == nil (mismatched types ST and nil)

}
