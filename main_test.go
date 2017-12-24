package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}

	a.Initialize("root", "sunday", "todo")

	ensureTableExist()

	var code = m.Run()

	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/todo", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestGetAllTodoList(t *testing.T) {
	clearTable()

	addListItem(5)

	req, _ := http.NewRequest("GET", "/api/todo", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetTodoListById(t *testing.T) {
	clearTable()

	addListItem(1)

	req, _ := http.NewRequest("GET", "/api/todo/1", nil)
	response := executeRequest(req)

	fmt.Println(response.Body)

	var s map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &s)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestAddTodoList(t *testing.T) {
	clearTable()

	todoListItem := []byte(`{"title": "blabla", "description": "bla bla bla", "done" : false}`)

	req, _ := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(todoListItem))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var item map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &item)

	if item["title"] != "blabla" {
		t.Errorf("failed added data title")
	}

	if item["description"] != "bla bla bla" {
		t.Errorf("failed added data description")
	}

	if item["done"] != false {
		t.Errorf("failed added data done")
	}
}

func TestAddEmptyTodoList(t *testing.T) {
	clearTable()

	todoListItem := []byte(`{"title": "", "description": "", "done" : }`)

	req, _ := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(todoListItem))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestUpdateTodoList(t *testing.T) {
	clearTable()

	addListItem(1)

	req, _ := http.NewRequest("GET", "/api/todo/1", nil)
	response := executeRequest(req)

	var originalUser map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &originalUser)

	updateItem := []byte(`{"title": "muehehe", "description": "muehehe", "done": true}`)

	req, _ = http.NewRequest("PUT", "/api/todo/1", bytes.NewBuffer(updateItem))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var updatedUser map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &updatedUser)

	if originalUser["id"] != updatedUser["id"] {
		t.Errorf("ID doesnt same. original : %v. Got %v", originalUser["id"], updatedUser["id"])
	}

	if originalUser["title"] == updatedUser["title"] {
		t.Errorf("title doesnt updated. Got : %v", updatedUser["title"])
	}

	if originalUser["description"] == updatedUser["description"] {
		t.Errorf("decription doesnt updated. Got : %v", updatedUser["description"])
	}

	if originalUser["done"] == updatedUser["done"] {
		t.Errorf("done doesnt updated. Got : %v", updatedUser["done"])
	}
}

func TestDeleteTodoList(t *testing.T) {
	clearTable()

	addListItem(1)

	req, _ := http.NewRequest("DELETE", "/api/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/todo/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}

func executeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, r)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d.", expected, actual)
	}
}

func addListItem(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO todolist(title, description, done) VALUES('%s', '%s', %d)", ("User" +
			strconv.Itoa(i+1)), "hehe", 0)
		a.DB.Exec(statement)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM todolist")
	a.DB.Exec("ALTER TABLE todolist AUTO_INCREMENT = 1")
}

func ensureTableExist() {
	_, err := a.DB.Exec(tableCreationQuery)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS todolist(
	id INT AUTO_INCREMENT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	done tinyint(1)	NOT NULL)`
