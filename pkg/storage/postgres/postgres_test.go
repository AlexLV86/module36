package postgres

import (
	"log"
	"module36/GoNews/pkg/storage"
	"os"
	"testing"
)

var s *Storage

func TestMain(m *testing.M) {
	path := os.Getenv("dbpath")
	if path == "" {
		m.Run()
	}
	var err error
	s, err = New(path)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
func TestStorage_Posts(t *testing.T) {
	post := []storage.Post{{Title: "Тестовая статья", Content: "Прекрасная тестовая статья для постгрес!",
		PubTime: 123242, Link: "https://ya3.ru"},
		{Title: "Тестовая статья", Content: "Прекрасная тестовая статья для постгрес!",
			PubTime: 123244524, Link: "https://ya4.ru"}}
	_, err := s.AddPosts(post)
	if err != nil {
		t.Fatal(err)
	}
	data, err := s.Posts(4)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestStorage_AddPosts(t *testing.T) {
	type args struct {
		posts []storage.Post
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AddPosts(tt.args.posts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.AddPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.AddPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}
