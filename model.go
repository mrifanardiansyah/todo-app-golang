package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ListItemPage struct {
	Title    string
	Header   string
	TodoList []ListItem
}

type ListItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	Done        bool   `json:"done"`
}

func (li *ListItem) GetList(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, title, description, done FROM todolist where id=%d;", li.ID)
	return db.QueryRow(statement).Scan(&li.ID, &li.Title, &li.Description, &li.Done)
}

func (li *ListItem) AddList(db *sql.DB) error {

	if li.Title == "" || li.Description == "" {
		return errors.New("Data tidak boleh kosong")
	}

	statement := fmt.Sprintf("INSERT INTO todolist (title, description, done) values ('%s', '%s', %t)", li.Title, li.Description, li.Done)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&li.ID)
	if err != nil {
		return err
	}

	return nil
}

func (li *ListItem) UpdateList(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE todolist set title='%s', description='%s', done=%t where id=%d",
		li.Title, li.Description, li.Done, li.ID)
	_, err := db.Exec(statement)
	return err
}

func GetAllList(db *sql.DB) ([]ListItem, error) {
	statement := fmt.Sprintf("select id, title, description, done from todolist;")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	List := []ListItem{}

	for rows.Next() {
		var i ListItem
		if err = rows.Scan(&i.ID, &i.Title, &i.Description, &i.Done); err != nil {
			return nil, err
		}
		List = append(List, i)
	}

	return List, nil
}

func (li *ListItem) DeleteList(db *sql.DB) error {
	statement := fmt.Sprintf("Delete from todolist where id=%d", li.ID)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}
