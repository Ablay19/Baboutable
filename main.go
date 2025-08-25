package main

import (
        "database/sql"
        "html/template"
        "log"
        "net/http"

        _ "modernc.org/sqlite"
)

type Person struct {
        ID    int
        Name  string
        Phone string
        Date  string
        Done  bool
}

var db *sql.DB
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
        var err error
        db, err = sql.Open("sqlite", "people.db")
        if err != nil {
                log.Fatal(err)
        }
        // إنشاء الجدول إذا ما كان موجود
        _, err = db.Exec(`CREATE TABLE IF NOT EXISTS people (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT,
                phone TEXT,
                date TEXT,
                done BOOLEAN
        )`)
        if err != nil {
                log.Fatal(err)
        }

        http.HandleFunc("/", listHandler)
        http.HandleFunc("/add", addHandler)
        http.HandleFunc("/delete", deleteHandler)
        http.HandleFunc("/toggle", toggleHandler)

        log.Println("🚀 الخادم يعمل على http://localhost:8080")
        http.ListenAndServe(":8080", nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
        rows, _ := db.Query("SELECT id, name, phone, date, done FROM people")
        defer rows.Close()

        var people []Person
        for rows.Next() {
                var p Person
                rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Date, &p.Done)
                people = append(people, p)
        }

        tmpl.Execute(w, people)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
                name := r.FormValue("name")
                phone := r.FormValue("phone")
                date := r.FormValue("date")
                done := r.FormValue("done") == "on"

                _, err := db.Exec("INSERT INTO people (name, phone, date, done) VALUES (?, ?, ?, ?)",
                        name, phone, date, done)
                if err != nil {
                        log.Println(err)
                }
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Query().Get("id")
        _, err := db.Exec("DELETE FROM people WHERE id = ?", id)
        if err != nil {
                log.Println(err)
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Query().Get("id")
        _, err := db.Exec("UPDATE people SET done = NOT done WHERE id = ?", id)
        if err != nil {
                log.Println(err)
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
}
