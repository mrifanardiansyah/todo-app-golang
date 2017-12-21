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
	Title       string
	Description string
	ID          int
	Done        bool
}

func (li *ListItem) GetList(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, title, description, done FROM todolist where id=%d;", li.ID)
	return db.QueryRow(statement).Scan(&li.ID, &li.Title, &li.Description, &li.Done)
}

func (li *ListItem) AddList(db *sql.DB) error {
	return errors.New("error. belum ditambah")
}

func (li *ListItem) UpdateList(db *sql.DB) error {
	return errors.New("error. belum ditambah")
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
	return errors.New("error. belum ditambah")
}
