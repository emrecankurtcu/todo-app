package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../html/index.html")
	})

	http.HandleFunc("/get-todo", getToDo)

	http.HandleFunc("/add-todo", addToDo)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":9990", nil); err != nil {
		log.Fatal(err)
	}
}

func addToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var data Todo
		json.Unmarshal(body, &data)

		response := addToDoFunc(data.Text)

		if response {
			w.Write([]byte(fmt.Sprintf("Added %s successfully!", data.Text)))
		} else {
			w.Write([]byte(fmt.Sprintf("Error!")))
		}

	}
}

func addToDoFunc(data string) bool {
	db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/todo")
	if err != nil {
		panic(err)
	}
	q, err := db.Prepare("INSERT INTO todo (text,is_done) VALUES (?,?)")
	if err != nil {
		panic(err)
	}

	_, err = q.Exec(data, "0")

	if err != nil {
		panic(err)
	}

	return true
}

func getToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/todo")
	if err != nil {
		panic(err)
	}
	var places []Todo
	rows, err := db.Query("SELECT text,is_done FROM todo")

	for rows.Next() {
		var place Todo = Todo{}
		if err := rows.Scan(&place.Text, &place.IsDone); err != nil {
			panic(err)
		}
		places = append(places, place)
	}

	rows.Scan(places)
	if err != nil {
		panic(err)
	}

	respBody, err := json.Marshal(places)

	if err != nil {
		panic(err)
	}

	w.Write(respBody)
}

type Todo struct {
	Text   string `db:"text"`
	IsDone bool   `db:"is_done"`
}
