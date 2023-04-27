package controllers
 
import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = ""
	password = ""
	dbname   = ""
)

type Book struct {
	BookID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Desc int `json:"desc"`
}

func getConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}


func CreateBook(ctx *gin.Context){
	var newBook Book;

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db, err := getConnection()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()
	sqlStatement := `INSERT INTO books (title, author, description) VALUES ($1, $2, $3) RETURNING id`
	var bookID int
	err = db.QueryRow(sqlStatement, newBook.Title, newBook.Author, newBook.Desc).Scan(&bookID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, "Created")
}

func UpdateBook(ctx *gin.Context){
	bookID := ctx.Param("bookID")
	var updatedBook Book

	if err := ctx.ShouldBindJSON(&updatedBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db, err := getConnection()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	sqlStatement := `UPDATE books SET title=$1, author=$2, description=$3 WHERE id=$4`
	res, err := db.Exec(sqlStatement, updatedBook.Title, updatedBook.Author, updatedBook.Desc, bookID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Data Not Found",
			"error_message": fmt.Sprintf("book with id %v not found", bookID),
		})
		return
	}

	ctx.JSON(http.StatusOK, "Updated")
}

func GetBook(ctx *gin.Context){
	bookID := ctx.Param("bookID")
	var bookData Book
	db, err := getConnection()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM books WHERE id=$1", bookID).Scan(&bookData.BookID, &bookData.Title, &bookData.Author, &bookData.Desc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status": "Data Not Found",
			"error_message": fmt.Sprintf("book with id %v not found", bookID),
		})
		return
	}

	ctx.JSON(http.StatusOK, bookData)
}

func GetBooks(ctx *gin.Context){
	db, err := getConnection()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Desc)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func DeleteBook(ctx *gin.Context){
	bookID := ctx.Param("bookID")
	
	db, err := getConnection()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	sqlStatement := "DELETE FROM books WHERE id=$1;"
	res, err := db.Exec(sqlStatement,bookID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Data Not Found",
			"error_message": fmt.Sprintf("book with id %v not found", bookID),
		})
		return
	}

	ctx.JSON(http.StatusOK, "Deleted")

}


