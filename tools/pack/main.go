package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/yinqiang/go-pizzaknife/knife"
)

func main() {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println("Error:", e.(error).Error())
		}
	}()

	infoPath := flag.String("f", "", "info file")
	outputPath := flag.String("o", "", "output file")
	flag.Parse()

	*infoPath = strings.TrimSpace(*infoPath)
	if len(*infoPath) == 0 {
		flag.Usage()
		return
	}
	*outputPath = strings.TrimSpace(*outputPath)

	knife.Pack(*infoPath, *outputPath)
}
