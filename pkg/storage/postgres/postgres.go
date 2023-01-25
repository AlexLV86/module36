package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"

	"module36/GoNews/pkg/storage"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Posts получение n публикаций
func (s *Storage) Posts(n int) ([]storage.Post, error) {
	query := `SELECT posts.id, posts.title, 
	posts.content, posts.pubdate, posts.link 
	FROM posts ORDER BY posts.pubdate DESC  LIMIT $1;`
	rows, err := s.db.Query(context.Background(), query, n)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

// AddPosts добавляет новые публикации
func (s *Storage) AddPosts(posts []storage.Post) (int, error) {
	var err error
	var cmd pgconn.CommandTag
	rows := 0
	for _, p := range posts {
		cmd, err = s.db.Exec(context.Background(), `
		INSERT INTO posts (title, content, pubdate, link) 
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;
		`,
			p.Title, p.Content, p.PubTime, p.Link)
		if err != nil {
			return 0, err
		}
		rows += int(cmd.RowsAffected())
	}
	return rows, err
}
