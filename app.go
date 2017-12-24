package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

func (a *App) Initialize(user, password, dbName string) {
	var connectionString = fmt.Sprintf("%s:%s@/%s", user, password, dbName)

	var err error

	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("cant open database")
	}

	a.Router = mux.NewRouter()
	a.InitializeRoute()
}

func (a *App) InitializeRoute() {
	a.Router.HandleFunc("/", a.redirect).Methods("GET")
	a.Router.HandleFunc("/api/todo", a.getTodoList).Methods("GET")
	a.Router.HandleFunc("/api/todo/{id:[0-9]+}", a.GetTodoListById).Methods("GET")
	a.Router.HandleFunc("/api/todo/{id:[0-9]+}", a.updateTodoList).Methods("PUT")
	a.Router.HandleFunc("/api/todo", a.addTodoList).Methods("POST")
	a.Router.HandleFunc("/api/todo/{id:[0-9]+}", a.deleteTodoList).Methods("DELETE")

	a.Router.HandleFunc("/todo", a.homePage).Methods("GET")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.itemPage).Methods("GET")
	a.Router.HandleFunc("/todo", a.addList).Methods("POST")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.updateList).Methods("POST")
	a.Router.HandleFunc("/todo/new", a.addPage).Methods("GET")
	a.Router.HandleFunc("/todo/edit/{id:[0-9]+}", a.editPage).Methods("GET")

	a.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (a *App) getTodoList(w http.ResponseWriter, r *http.Request) {
	List, err := GetAllList(a.DB)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJson(w, http.StatusOK, List)
}

func (a *App) GetTodoListById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid Id TodoListItem")
	}

	var item ListItem
	item.ID = id
	err = item.GetList(a.DB)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
	}

	responseWithJson(w, http.StatusOK, item)
}

func (a *App) updateTodoList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var li ListItem

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&li); err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	li.ID = id

	if err := li.UpdateList(a.DB); err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJson(w, http.StatusOK, li)
}

func (a *App) addTodoList(w http.ResponseWriter, r *http.Request) {
	var li ListItem

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&li); err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := li.AddList(a.DB); err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJson(w, http.StatusOK, li)
}

func (a *App) deleteTodoList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var li ListItem
	li.ID = id

	if err = li.DeleteList(a.DB); err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJson(w, http.StatusOK, nil)
}

func (a *App) homePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		responseWithError(w, http.StatusNotFound, err.Error())
		return
	}

	list, _ := GetAllList(a.DB)
	t.Execute(w, ListItemPage{Header: "Daftar Tugas dan Kuis", Title: "Daftar Tugas dan Kuis", TodoList: list})
}

func (a *App) itemPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	t, err := template.ParseFiles("./templates/item-description.html")
	if err != nil {
		responseWithError(w, http.StatusNotFound, err.Error())
		return
	}

	var item ListItem
	item.ID = id
	err = item.GetList(a.DB)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t.Execute(w, item)
}

func (a *App) addList(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	payload := []byte(fmt.Sprintf(`{"title" : "%s", "description" : "%s", "done" : %t}`, title, description, false))

	req, err := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(payload))
	a.executeRequest(req)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (a *App) updateList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := r.ParseForm(); err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	payload := []byte(fmt.Sprintf(`{"title" : "%s", "description" : "%s", "done" : %t}`, title, description, false))

	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/todo/%d", id), bytes.NewBuffer(payload))
	a.executeRequest(req)

	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (a *App) addPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/addlistitem.html")
	if err != nil {
		responseWithError(w, http.StatusNotFound, err.Error())
		return
	}
	t.Execute(w, map[string]string{"Header": "Add item to list"})
}

func (a *App) editPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var li ListItem
	li.ID = id

	if err = li.GetList(a.DB); err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t, err := template.ParseFiles("./templates/editlistitem.html")
	if err != nil {
		responseWithError(w, http.StatusNotFound, err.Error())
		return
	}
	t.Execute(w, map[string]string{"Header": "Edit item to list", "Title": li.Title,
		"Description": li.Description, "ID": strconv.Itoa(id)})
}

func responseWithError(w http.ResponseWriter, code int, message string) {
	responseWithJson(w, code, map[string]string{"error": message})
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) executeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, r)

	return rr
}
