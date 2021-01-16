package mysql

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"naz.net/snippet/pkg/models"
	"strconv"
	"time"
)

type SnippetModel struct {
	DB *pgxpool.Pool
}

const (
	insertSql = "INSERT INTO snippets (title, content, created, expires) VALUES($1, $2, $3,$4) RETURNING id"
	select1   = "select id,title, content, created, expires from snippets where id=$1"
	select2   = "SELECT id,title, content, created, expires from snippets order by created DESC LIMIT 10"
)

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	var id uint64
	days, err := strconv.Atoi(expires)
	if err != nil {
		return 0, err
	}
	row := m.DB.QueryRow(context.Background(), insertSql, title, content, time.Now(), time.Now().AddDate(0, 0, days))
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	err := m.DB.QueryRow(context.Background(), select1, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	snippets := []*models.Snippet{}
	rows, err := m.DB.Query(context.Background(), select2)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
