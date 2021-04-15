package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age  uint16 `json:"age"`
}

var users = []User{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "mysql:mysql@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("select `name`, `age` from `users`")
	if err != nil {
		panic(err)
	}

	users = nil

	for res.Next() {
		var user User
		err = res.Scan(&user.Name, &user.Age)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	t.ExecuteTemplate(w, "index", users)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	age := r.FormValue("age")

	if name == "" || age == "" {
		fmt.Fprintf(w, "Поле не заполнено")
	} else {
		db, err := sql.Open("mysql", "mysql:mysql@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("insert into `users` (`name`, `age`) values ('%s', '%s')", name, age))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleFunc() {
	go http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	go http.HandleFunc("/save", save)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
