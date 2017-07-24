// Copyright 2017 blank. all rights reserved.
// Authors: blank

/*
Package film 提供了一个爬电影的逻辑，根据不同的网站，爬取数据并存储
*/
package film

import (
	"net/url"
	"strings"
	"sync"
)

var filmMutex sync.Mutex
var films map[string]Film

// FilterFunc 定义了在爬取电影的同，过滤掉返回false数据
type FilterFunc func(data Data) bool

// FilmType 电影类型 动作
type FilmType uint32

const (
	// Action 动作类型
	Action FilmType = 1 << iota
)

// FilmLanguage 电影语言
type FilmLanguage uint32

const (
	// English 语种为英语
	English FilmLanguage = 1 << iota
)

// Option 爬取电影的选项
type Option struct {
	LimitCount uint32       // 最多获取多少条数据
	Offset     uint32       // 获取数据的偏移位置
	Type       FilmType     // 电影类型
	Language   FilmLanguage // 电影语种
}

// Film 定义取电影的接口
type Film interface {
	Spider(filmURL string, option Option, filters ...FilterFunc) <-chan Data
}

// Data 定义了存储电影的接口
type Data interface {
	DownloadAddr() DownloadAddrs
	Title() string
	Score() float64
	ImageURL() string
	EM() string
}

// DownloadAddrs 定义下载地址集合
type DownloadAddrs interface {
	All() map[string]string
	Each(func(name, URL string))
}

// Spider 提供外部一个接口用来获取对应网站的数据
func Spider(filmURL string, option Option) <-chan Data {
	if f := getFilm(filmURL); f != nil {
		return f.Spider(filmURL, option)
	}
	return nil
}

func getFilm(filmURL string) Film {
	u, err := url.Parse(filmURL)
	if err != nil {
		return nil
	}
	h := u.Host
	if i := strings.LastIndex(h, ":"); i != -1 {
		h = h[:i]
	}
	if i := strings.LastIndex(h, `/`); i != -1 {
		h = h[i:]
	}

	hs := strings.Split(h, `.`)
	for _, s := range hs {
		if f, ok := films[s]; ok {
			return f
		}
	}
	return nil
}

func registerSpider(name string, f Film) {
	filmMutex.Lock()
	defer filmMutex.Unlock()
	if films == nil {
		films = make(map[string]Film, 1)
	}
	films[name] = f
}
