// API приложения Agrigator.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	storage "Agrigator/pkg/storage/pstg"

	"github.com/gorilla/mux"
)

type API struct {
	db *storage.DB
	r  *mux.Router
}

func New(db *storage.DB) *API {
	a := API{db: db, r: mux.NewRouter()}
	a.endpoints()
	return &a

}

func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	api.r.HandleFunc("/news/last", api.lastArticles).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/lastlist", api.lastArticlesList).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/filter", api.newsFilteredList).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/news", api.newsFullDetailed).Methods(http.MethodGet, http.MethodOptions)
}

// lastArticles обрабатывает запрос на последние n новостгй с деталями
// в браузере http://localhost:998/news/last?n=5
// n
func (api *API) lastArticles(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query()
	ns := form.Get("n")
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	news, err := api.db.LastArticles(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
	log.Println("Agrigator: API:lastArticles:", "ok ", r.URL.Query().Encode())
}

// lastArticlesList обрабатывает запрос на последние n новостгй с деталями
// в браузере http://localhost:998/news/lastlist?n=5
// n
func (api *API) lastArticlesList(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query()
	ns := form.Get("n")
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	news, err := api.db.LastArticlesList(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
	log.Println("Agrigator: API:lastArticlesList:", "ok ", r.URL.Query().Encode())
}

// newsFilteredDetailed обрабатывает запрос на поиск новости по полям title или
// в браузере http://localhost:998/news/filter?time1=1699016144&time2=1700293140&lim=100&field=title&contains=putin&sortfield=id&dir=s
// time1, time2, Lim, field, content, sortfield, dir
func (api *API) newsFilteredList(w http.ResponseWriter, r *http.Request) {
	var fParam storage.FilterParam
	var err error
	form := r.URL.Query()
	fParam.Time[0] = form.Get("time1")
	fParam.Time[1] = form.Get("time2")
	fParam.N = form.Get("lim")
	fParam.Field = form.Get("field")
	fParam.Contains = "'%" + form.Get("contains") + "%'"
	fParam.Sort.Field = form.Get("sortfield")
	dirs := form.Get("dir")
	fParam.Sort.DirUp = ""
	if dirs != "" {
		fParam.Sort.DirUp = "DESC"
	}

	news, err := api.db.NewsFilteredlist(fParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
	log.Println("Agrigator: API:newsFilteredList:", "ok ", r.URL.Query().Encode())
}

// NewsFullDetailed обрабатывает запрос на конкретную новость с деталями
// в браузере http://localhost:998/news/news?id=5
// id
func (api *API) newsFullDetailed(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query()
	ids := form.Get("id")
	n, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	news, err := api.db.NewsFullDetailed(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
	log.Println("Agrigator: API:newsFullDetailed:", "ok ", r.URL.Query().Encode())
}
