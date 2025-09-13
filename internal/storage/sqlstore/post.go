package sqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RomanKovalev007/mai_news/internal/models"
	"github.com/RomanKovalev007/mai_news/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func (s *Storage) GetAllPosts() ([]models.OutputPost, error){
	op := "storage.sqlstore.GetAllPosts"

	stmt, err := s.db.Prepare("SELECT id, title, content, created_at FROM post")
		if err != nil {
			return []models.OutputPost{}, fmt.Errorf("%s: prepare statement: %w", op, err)
		}

	rows, err := stmt.Query()
	if err != nil{
		return []models.OutputPost{}, fmt.Errorf("%s: failed to get all posts: %w", op, err)
	}

	var posts []models.OutputPost

	for rows.Next(){
		var post models.OutputPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return []models.OutputPost{}, fmt.Errorf("%s: scan row: %w", op, err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil{
		return []models.OutputPost{}, fmt.Errorf("%s: rows err: %w", op, err)
	}
	
	return posts, nil
}

func (s *Storage) SavePost(inputPost models.InputPost) (models.OutputPost, error){
	op := "storage.sqlstore.SavePost"

	stmt, err := s.db.Prepare("INSERT INTO post(title, content, created_at) VALUES(?, ?, ?)")
	if err != nil {
		return models.OutputPost{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	now := time.Now()
	res, err := stmt.Exec(inputPost.Title, inputPost.Content, now)
	if err != nil {
		return models.OutputPost{}, fmt.Errorf("%s: exec statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.OutputPost{}, fmt.Errorf("%s: get last insert id: %w", op, err)
	}

	post := models.OutputPost{
		ID: int(id),
		Title: inputPost.Title,
		Content: inputPost.Content,
		CreatedAt: now.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
	}

	return post, nil
}

func (s *Storage) GetPost(id int) (models.OutputPost, error){
	op := "storage.sqlstore.GetPost"

	stmt, err := s.db.Prepare("SELECT id, title, content, created_at FROM post WHERE id = ?")
	if err != nil{
		return models.OutputPost{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	row := stmt.QueryRow(id)
	var post models.OutputPost
	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return models.OutputPost{}, storage.ErrPostNotFound
		}
		return models.OutputPost{}, fmt.Errorf("%s: scan row: %w", op, err)
	}

	return post, nil
}

func (s *Storage) PatchPost(id int, inputPost models.InputPost) (models.OutputPost, error){
	op := "storage.sqlstore.PutchPost"

	stmt, err := s.db.Prepare(`
	UPDATE post SET title = ?, content = ? WHERE id = ?
	RETURNING id, title, content, created_at`)
	if err != nil{
		return models.OutputPost{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var post models.OutputPost
	err = stmt.QueryRow(inputPost.Title, inputPost.Content, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return models.OutputPost{}, storage.ErrPostNotFound
		}
		return models.OutputPost{}, fmt.Errorf("%s: scan row: %w", op, err)
	}

	return post, nil
}

func (s *Storage) DeletePost(id int) error{
	op := "storage.sqlstore.DeletePost"

	stmt, err := s.db.Prepare("DELETE FROM post WHERE id = ?")
	if err != nil{
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return storage.ErrPostNotFound
		}
		return fmt.Errorf("%s: failed delete: %w", op, err)
	}

	return nil
}