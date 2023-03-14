package main

import (
    "fmt"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
}

func main() {
	db, err := sql.Open("mysql", "root:Revanth@1436@tcp(172.17.0.5:3306)/Books")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", getBook(db)).Methods("GET")
	router.HandleFunc("/books", createBook(db)).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook(db)).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook(db)).Methods("DELETE")
    fmt.Println("connected")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM books")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var books []Book

		for rows.Next() {
			var book Book
			if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
				log.Fatal(err)
			}
			books = append(books, book)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(books)
	}
}

func getBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		result := db.QueryRow("SELECT id, title, author FROM books WHERE id = ?", params["id"])

		var book Book

		err := result.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(book)
	}
}

func createBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book Book
		json.NewDecoder(r.Body).Decode(&book)

		_, err := db.Exec("INSERT INTO books(id, title, author) VALUES(?,?,?)", book.ID, book.Title, book.Author)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(book)
	}
}

func updateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var book Book
		json.NewDecoder(r.Body).Decode(&book)

		_, err := db.Exec("UPDATE books SET title = ?, author = ? WHERE id = ?", book.Title, book.Author, params["id"])
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(book)
	}
}

func deleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		_, err := db.Exec("DELETE FROM books WHERE id = ?", params["id"])
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode("Book deleted successfully")
	}
}
