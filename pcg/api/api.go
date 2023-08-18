package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"GoNews/pcg/database"

	"github.com/gorilla/mux"
)

type API struct {
	r        *mux.Router
	db       *sql.DB
	rssLinks []string
}

func NewAPI(db *sql.DB, rssLinks []string) *API {
	api := &API{
		r:        mux.NewRouter(),
		db:       db,
		rssLinks: rssLinks,
	}

	api.endpoints()
	return api
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.r.ServeHTTP(w, r)
}

func (api *API) GetRouter() *mux.Router {
	return api.r
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		http.Error(w, "Неверное количество новостей", http.StatusBadRequest)
		return
	}

	// Получить n последних новостей из базы данных
	posts, err := database.GetLatestPosts(n)
	if err != nil {
		http.Error(w, "Не удалось получить новости", http.StatusInternalServerError)
		return
	}

	// Отправить ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (api *API) webAppHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./webapp")).ServeHTTP(w, r)
}

func (api *API) endpoints() {
	api.r.HandleFunc("/news/{n:[0-9]+}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	api.r.PathPrefix("/").HandlerFunc(api.webAppHandler).Methods(http.MethodGet)
}

func StartAPI(port string, db *sql.DB, rssLinks []string) error {
	api := NewAPI(db, rssLinks)
	return http.ListenAndServe(":"+port, api)
}