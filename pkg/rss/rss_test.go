package rss

import "testing"

func TestRSS_Get(t *testing.T) {
	r := New()
	url := "https://habr.com/ru/rss/best/daily/?fl=ru"
	err := r.Get(url)
	if err != nil {
		t.Fatalf("Ошибка при получении rss: %v", err)
	}
	var got, nowant int64 = 0, 0
	got = int64(len(r.Channel.Item))
	if got == nowant {
		t.Fatalf("got %d, no want %d", got, nowant)
	}
	pubDate := "Mon, 2 Jan 2006 15:04:05 -0700"
	got, err = r.DateToUnix(pubDate)
	if err != nil {
		t.Fatalf("Ошибка конвертации даты: %v", err)
	}
	want := int64(1136239445)
	if got != want {
		t.Fatalf("got %d, no want %d", got, want)
	}
}
