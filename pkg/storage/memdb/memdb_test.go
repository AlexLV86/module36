package memdb

import (
	"module36/GoNews/pkg/storage"
	"testing"
)

func TestDB_Posts(t *testing.T) {
	db := New()
	p := []storage.Post{{Title: "Test title", Content: "Test content", PubTime: 123214},
		{Title: "Test title 2", Content: "Test content 2", PubTime: 123214}}
	got, err := db.AddPosts(p)
	want := 2
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	posts, err := db.Posts(10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(posts)

	posts, err = db.Posts(0)
	if err != nil {
		t.Fatal(err)
	}
	got = len(posts)
	want = 0
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
