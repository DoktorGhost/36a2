package main

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

var db *sql.DB

func initDB() {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)
	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	schemaSQL, err := os.ReadFile("schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		log.Fatal(err)
	}

	// Добавьте ограничение уникальности на поле title
	_, err = db.Exec("ALTER TABLE news ADD CONSTRAINT unique_title UNIQUE (title)")
	if err != nil {
		log.Fatal(err)
	}
}
func saveToDB(post typeStruct.Post) error {
	query := `
		INSERT INTO news (title, description, pub_date, source)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	row := db.QueryRow(query, post.Title, post.Content, post.PubTime, post.Link)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	fmt.Println("Saved to DB with ID:", id)
	return nil
}

func readFromDB(title string) (typeStruct.Post, error) {
	var post typeStruct.Post

	query := `
		SELECT id, title, description, pub_date, source
		FROM news
		WHERE title = $1
	`
	row := db.QueryRow(query, title)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return post, err
	}

	return post, nil
}
