package main

import (
	"fmt"
	"moive/film"
)

func main() {
	opt := film.Option{
		LimitCount: 24,
		Offset:     0,
		Language:   film.English,
		Type:       film.Action,
	}
	datas := film.Spider("https://m.80s.tw/", opt)
	for data := range datas {
		fmt.Println(data)
	}
}
