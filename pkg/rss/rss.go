package rss

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ItemStruct struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Content struct {
	XMLName xml.Name     `xml:"channel"`
	Title   string       `xml:"title"`
	Item    []ItemStruct `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Content  `xml:"channel"`
}

func New() *RSS {
	return &RSS{}
}

// получаем публикации из rss канала по url
func (r *RSS) Get(url string) error {
	// подготавливаем запрос на получение rss
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	// получаем ответ от rss канала
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// читаем тело ответа в xml формате
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//var articles RSS
	err = xml.Unmarshal(text, r)
	if err != nil {
		return err
	}
	return nil
}

// конвертер даты из формата
// Mon, 2 Jan 2006 15:04:05 -0700
// в UNIX
func (r *RSS) DateToUnix(pubDate string) (int64, error) {
	layoutUTC := "Mon, 2 Jan 2006 15:04:05 -0700"
	layoutGMT := "Mon, 2 Jan 2006 15:04:05 GMT"
	t1p, err := time.Parse(layoutUTC, pubDate)
	if err != nil {
		if t1p, err = time.Parse(layoutGMT, pubDate); err != nil {
			return 0, fmt.Errorf("parse date error %s, ", pubDate)
		}
	}
	return t1p.Unix(), nil
}
