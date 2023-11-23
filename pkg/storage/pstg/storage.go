// работа с БД PostgreSQL
package storage

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// База данных.
type DB struct {
	pool *pgxpool.Pool
}

// тип сортировки
type SortParam struct {
	Field string
	DirUp string // направление DESC или ""
}

// параметр фильтрации
type FilterParam struct {
	Time     [2]string // период времени
	N        string    // число в списке
	Field    string    // поле для поиска
	Contains string    // содержит...
	Sort     SortParam // параметры сортировки

}

// Публикация, получаемая из RSS.
type Article struct {
	ID        int    // номер записи
	Title     string // заголовок публикации
	Content   string // содержание публикации
	PubTime   int64  // время публикации
	Url       string // ссылка на источник
	Publisher string // название источника
	Autor     string // Имя автора
}

func New(adr string) (*DB, error) {
	dbuserpass := os.Getenv("agrigatordb") //"postgres://postgres:" + "password"
	connstr := dbuserpass + adr
	if connstr == "" {
		return nil, errors.New("не указано подключение к БД")
	}
	pool, err := pgxpool.Connect(context.Background(), connstr)

	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

// SaveArticles Сохраниеть статью в БД
func (db *DB) SaveArticles(articles []Article) error {
	for _, article := range articles {
		_, err := db.pool.Exec(context.Background(), `
		INSERT INTO articles (title, content, pub_time, a_url, publisher, author)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (title) DO NOTHING`,
			article.Title,
			article.Content,
			article.PubTime,
			article.Url,
			article.Publisher,
			article.Autor,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// LastArticles возвращает последние N новостей со всеми данными из БД.
func (db *DB) LastArticles(n int) ([]Article, error) {
	if n == 0 {
		n = 10
	}
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, content, pub_time, a_url, publisher, author FROM articles
	ORDER BY pub_time DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}
	var alist []Article
	for rows.Next() {
		var p Article
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Url,
			&p.Publisher,
			&p.Autor,
		)
		if err != nil {
			return nil, err
		}
		alist = append(alist, p)
	}
	return alist, rows.Err()
}

// LastArticlesList возвращает список последних N тизеров новостей из БД
// (ID, заголовок и время публикации).
func (db *DB) LastArticlesList(n int) ([]Article, error) {
	if n == 0 {
		n = 9223372036854775807
	}
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, pub_time FROM articles
	ORDER BY pub_time DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}
	var alist []Article
	for rows.Next() {
		var p Article
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.PubTime,
		)
		if err != nil {
			return nil, err
		}
		alist = append(alist, p)
	}
	return alist, rows.Err()
}

// NewsFilteredlist возвращает отфильтрованный и отсортированный список новостей
func (db *DB) NewsFilteredlist(f FilterParam) ([]Article, error) {
	if f.N == "0" {
		f.N = "99999"
	}
	qstr := "SELECT id, title, pub_time FROM articles " +
		"WHERE pub_time BETWEEN " + f.Time[0] + " AND " + f.Time[1] +
		" AND lower(" + f.Field + ") LIKE lower(" + f.Contains + ")" +
		" ORDER BY " + f.Sort.Field + " " + f.Sort.DirUp +
		" LIMIT " + f.N

	rows, err := db.pool.Query(context.Background(), qstr)
	if err != nil {
		return nil, err
	}
	var alist []Article
	for rows.Next() {
		var p Article
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.PubTime,
		)
		if err != nil {
			return nil, err
		}
		alist = append(alist, p)
	}
	return alist, rows.Err()
}

// NewsFullDetailed возвращает конкретную новость
func (db *DB) NewsFullDetailed(id int) (Article, error) {
	var a Article
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, content, pub_time, a_url, publisher, author FROM articles
	WHERE id = $1
	`,
		id,
	)
	if err != nil {
		return a, err
	}
	for rows.Next() {
		var p Article
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Url,
			&p.Publisher,
			&p.Autor,
		)
		if err != nil {
			return a, err
		}
		a = p
	}
	return a, rows.Err()
}
