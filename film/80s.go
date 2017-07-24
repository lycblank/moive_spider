package film

import (
	"bytes"
	"fmt"
	"sync"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	pageNum uint32 = 24
)

type film80sData struct {
	imgURL        string
	detailURL     string
	title         string
	em            string // 强调的内容
	downloadAddrs map[string]string
}

type film80sDownloadAddrs struct {
	downloadAddrs map[string]string
}

func (f80da *film80sDownloadAddrs) All() map[string]string {
	return f80da.downloadAddrs
}

func (f80da *film80sDownloadAddrs) Each(f func(name, URL string)) {
	for key, val := range f80da.downloadAddrs {
		f(key, val)
	}
}

func (f80sd *film80sData) DownloadAddr() DownloadAddrs {
	return &film80sDownloadAddrs{
		downloadAddrs: f80sd.downloadAddrs,
	}
}
func (f80sd *film80sData) ImageURL() string {
	return f80sd.imgURL
}
func (f80sd *film80sData) Title() string {
	return f80sd.title
}
func (f80sd *film80sData) Score() float64 {
	return 0.0
}
func (f80sd *film80sData) EM() string {
	return f80sd.em
}

type film80s struct {
}

func init() {
	registerSpider("80s", &film80s{})
}

func (f *film80s) Spider(filmURL string, option Option, filters ...FilterFunc) <-chan Data {
	if filmURL != "" && filmURL[len(filmURL)-1:] != `/` {
		filmURL = filmURL + `/`
	}
	listData := make(chan film80sData, 1024)
	// 列表数据解析
	wg := sync.WaitGroup{}
	wg.Add(int(option.LimitCount/pageNum + 1))
	listFunc := func(offset, limit uint32) {
		defer wg.Done()
		op := option
		op.Offset = offset
		op.LimitCount = limit
		listURL := f.buildURL(filmURL, op)
		f.parseList(filmURL, listURL, listData)
	}
	i := uint32(0)
	for ; i < option.LimitCount/pageNum+1; i++ {
		if (i+1)*pageNum > option.LimitCount {
			go listFunc(option.Offset+i*pageNum, option.LimitCount-i*pageNum)
		} else {
			go listFunc(option.Offset+i*pageNum, pageNum)
		}
	}
	go func() {
		wg.Wait()
		close(listData)
	}()
	// 详情数据解析
	return f.parseDetail(listData)
}

func (f *film80s) parseList(filmURL, listURL string, listData chan<- film80sData) {
	doc, err := goquery.NewDocument(listURL)
	if err != nil {
		fmt.Println(err, listURL)
		return
	}
	doc.Find(".col-xs-4").Each(func(i int, s *goquery.Selection) {
		data := film80sData{}
		s.Find("a").Each(func(i int, ss *goquery.Selection) {
			if href, exists := ss.Attr("href"); exists && href != "" {
				data.detailURL = fmt.Sprintf("%s%s", filmURL, href)
			}
			ss.Find("img").Each(func(i int, sss *goquery.Selection) {
				src, exists := sss.Attr("data-original")
				if exists && src != "" {
					if !strings.Contains(src, "https:") {
						src = fmt.Sprintf("%s%s", "https:", src)
					}
					data.imgURL = src
				}
			})
		})

		if data.detailURL == "" {
			return
		}
		s.Find(".list_mov_title").Each(func(i int, ss *goquery.Selection) {
			data.title = ss.Find("h4 a").Text()
			data.em = ss.Find("em").Text()
		})
		listData <- data
	})
}

func (f *film80s) parseDetail(listData <-chan film80sData) <-chan Data {
	datas := make(chan Data, 1024)
	go func() {
		defer close(datas)
		for ldata := range listData {
			data := ldata
			datas <- &data
		}
	}()
	return datas
}

func (f *film80s) buildURL(filmURL string, option Option) string {
	var buf bytes.Buffer
	buf.WriteString(filmURL)
	buf.WriteString("movie/")
	one, two, three, four, five, six, seven := 1, int(option.Offset), 0, 0, 0, 0, 0
	if option.Language == English {
		five = 2
	}
	if option.Type == Action {
		three = 1
	}
	buf.WriteString(fmt.Sprintf("%d-%d-%d-%d-%d-%d-%d", one, two, three, four, five, six, seven))
	return buf.String()
}
