package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"GoNews/pcg/typeStruct"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	db *sql.DB
}

func NewAPI(db *sql.DB) *API {
	r := mux.NewRouter()
	api := &API{r: r, db: db}

	// Регистрация методов API в маршрутизаторе запросов.
	api.endpoints()

	return api
}

func (api *API) endpoints() {
	api.r.HandleFunc("/news/{n}", api.getLatestNews).Methods(http.MethodGet, http.MethodOptions)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) getLatestNews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n := vars["n"]

	// Преобразование n в число (количество новостей)

	count, err := strconv.Atoi(n)
	if err != nil {
		log.Println(err)
	}

	// Запрос к базе данных для получения n последних новостей
	query := `
		SELECT title, description, pub_date, source
		FROM news
		ORDER BY pub_date DESC
		LIMIT $1
	`
	rows, err := api.db.Query(query, count)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	var news []typeStruct.Post
	for rows.Next() {
		var post typeStruct.Post
		err := rows.Scan(&post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		news = append(news, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func main() {
	db, err := sql.Open("postgres", "...")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api := NewAPI(db)

	serverAddr := "localhost:8080"
	fmt.Printf("Server is running at %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, api.r))
}
