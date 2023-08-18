package database

import (
	"GoNews/pcg/typeStruct"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	DBHost     = "localhost"
	DBPort     = "5432"
	DBUser     = "test_user"
	DBPassword = "qwerty123"
	DBName     = "testdb"
)

var DB *sql.DB

func InitDB() {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)
	var err error
	DB, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTable() {
	schemaSQL, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(string(schemaSQL))
	if err != nil {
		log.Fatal(err)
	}

	// Добавьте ограничение уникальности на поле title
	_, err = DB.Exec("ALTER TABLE news ADD CONSTRAINT unique_title UNIQUE (title)")
	if err != nil {
		log.Fatal(err)
	}
}

func SaveToDB(post typeStruct.Post) error {
	query := `
		INSERT INTO news (title, description, pub_date, source)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	row := DB.QueryRow(query, post.Title, post.Content, post.PubTime, post.Link)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	fmt.Println("Saved to DB with ID:", id)
	return nil
}

func ReadFromDB(title string) (typeStruct.Post, error) {
	var post typeStruct.Post

	query := `
		SELECT id, title, description, pub_date, source
		FROM news
		WHERE title = $1
	`
	row := DB.QueryRow(query, title)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return post, err
	}

	return post, nil
}

func GetLatestPosts(n int) ([]typeStruct.Post, error) {
	query := `
		SELECT id, title, description, pub_date, source
		FROM news
		ORDER BY pub_date DESC
		LIMIT $1
	`
	rows, err := DB.Query(query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []typeStruct.Post

	for rows.Next() {
		var post typeStruct.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func DeleteByTitle(title string) error {
	_, err := DB.Exec("DELETE FROM news WHERE title = $1", title)
	return err
}
