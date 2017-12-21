package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
	a.Router.HandleFunc("/todo", a.getTodoList).Methods("GET")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.GetTodoListById).Methods("GET")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.updateTodoList).Methods("PUT")
	a.Router.HandleFunc("/todo", a.addTodoList).Methods("POST")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.deleteTodoList).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getTodoList(w http.ResponseWriter, r *http.Request) {
	List, err := GetAllList(a.DB)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var todoListPage = ListItemPage{Header: "Todo List", Title: "Todo List", TodoList: List}

	t, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		responseWithError(w, http.StatusNotFound, err.Error())
		return
	}
	t.Execute(w, todoListPage)

	// responseWithJson(w, http.StatusOK, List)
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

}

func (a *App) addTodoList(w http.ResponseWriter, r *http.Request) {

}

func (a *App) deleteTodoList(w http.ResponseWriter, r *http.Request) {

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
