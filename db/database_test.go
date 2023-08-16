package main

import (
	"testing"
	"time"

	"GoNews/pcg/typeStruct"
)

func TestSaveAndReadFromDB(t *testing.T) {
	initDB()

	defer db.Close()

	// Создаем тестовый пост
	testPost := typeStruct.Post{
		Title:   "Test Title",
		Content: "Test Content",
		PubTime: time.Unix(16, 0),
		Link:    "http://example.com/test",
	}

	// Сохраняем тестовый пост в базу данных
	err := saveToDB(testPost)
	if err != nil {
		t.Fatalf("Failed to save post to DB: %v", err)
	}

	// Читаем пост из базы данных по названию
	readPost, err := readFromDB("Test Title") // Используем название для поиска
	if err != nil {
		t.Fatalf("Failed to read post from DB: %v", err)
	}

	// Сравниваем значения, исключая дату публикации
	if readPost.Title != testPost.Title ||
		readPost.Content != testPost.Content ||
		readPost.Link != testPost.Link {
		t.Errorf("Saved data doesn't match expected data.")
	}
}
