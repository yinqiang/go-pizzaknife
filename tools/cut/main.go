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

	fpath := flag.String("f", "", "file path")
	size := flag.Int64("s", 1024*1024*10, "part size")
	flag.Parse()

	*fpath = strings.TrimSpace(*fpath)
	if 0 == len(*fpath) {
		flag.Usage()
		return
	}

	knife.CutBySize(*fpath, *size)
}
