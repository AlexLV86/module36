package memdb

import (
	"module36/GoNews/pkg/storage"
	"sync"
)

// База данных заказов.
type DB struct {
	m    sync.Mutex           //мьютекс для синхронизации доступа
	id   int                  // текущее значение ID для нового заказа
	post map[int]storage.Post // БД заказов
}

// Конструктор БД.
func New() *DB {
	db := DB{
		id:   1, // первый номер заказа
		post: map[int]storage.Post{},
	}
	return &db
}

// Posts возвращает заданное кол-во статей
func (db *DB) Posts(n int) ([]storage.Post, error) {
	db.m.Lock()
	defer db.m.Unlock()
	data := make([]storage.Post, 0, n)
	i := 0
	for _, p := range db.post {
		if i >= n {
			break
		}
		i++
		// отправляю  n новостей без сортировки по дате
		data = append(data, p)
	}
	return data, nil
}

// AddPost добавляет статью без проверки на уникальность
func (db *DB) AddPosts(posts []storage.Post) (int, error) {
	db.m.Lock()
	defer db.m.Unlock()
	for _, p := range posts {
		p.ID = db.id
		db.post[p.ID] = p
		db.id++
	}
	return len(posts), nil
}

// размер БД
func (db *DB) Len() int {
	return len(db.post)
}
