package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"module36/GoNews/pkg/api"
	"module36/GoNews/pkg/rss"
	"module36/GoNews/pkg/storage"
	"module36/GoNews/pkg/storage/memdb"
)

type Config struct {
	ListRss   []string `json:"rss"`
	ReqPeriod int      `json:"request_period"`
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	db := memdb.New()
	chErr := make(chan error)
	chPosts := make(chan []storage.Post)
	// запускаем обработчик ошибок
	go handlerError(chErr)
	// запускаем обработчик для добавления новостей в БД
	go handlerDBPosts(chErr, chPosts, db)
	// читаем новости и засыпаем на время из config.json
	go func() {
		for {
			for _, url := range config.ListRss {
				go readRSS(chErr, chPosts, url)
			}
			time.Sleep(time.Minute * time.Duration(config.ReqPeriod))
		}
	}()

	// Создание объекта API, использующего БД в памяти.
	api := api.New(db)
	// Запуск сетевой службы и HTTP-сервера
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

// читаем RSS
func readRSS(chErr chan<- error, chPosts chan<- []storage.Post, url string) {
	r := rss.New()
	err := r.Get(url)
	if err != nil {
		// пишем ошибку в канал
		chErr <- err
		return
	}
	posts := make([]storage.Post, 0, len(r.Channel.Item))
	var p storage.Post
	// добавляем каждую новость в storage.Post
	for _, v := range r.Channel.Item {
		p.PubTime, err = r.DateToUnix(v.PubDate)
		if err != nil {
			// передать в канал ошибку и continue
			chErr <- fmt.Errorf("date convert error %s, wrap error: %w", v.PubDate, err)
			continue
		}
		p.Content = v.Description
		p.Link = v.Link
		p.Title = v.Title
		posts = append(posts, p)
	}
	// передаем все новости в канал
	chPosts <- posts
}

// Обработчик новостей. Добавляем новости в БД
func handlerDBPosts(chErr chan<- error, chPosts <-chan []storage.Post, db storage.Interface) {
	for posts := range chPosts {
		_, err := db.AddPosts(posts)
		if err != nil {
			// передать в канал ошибку и continue
			chErr <- err
		}
	}
}

// Обработчик ошибок. Записываем все ошибки в файл
func handlerError(chErr <-chan error) {
	for e := range chErr {
		f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE, os.ModePerm) //os.O_WRONLY|os.O_CREATE
		if err != nil {
			panic(err)
		}
		str := fmt.Sprintf("%v\n", e)
		f.WriteString(str)
	}
}

// читаем файл с конфигурацией rss каналов
func readConfig() (Config, error) {
	// достать данные из файла для
	c := Config{}
	data, err := ioutil.ReadFile("cmd/gonews/config.json")
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		fmt.Println(err)
		return Config{}, err
	}
	return c, nil
}
