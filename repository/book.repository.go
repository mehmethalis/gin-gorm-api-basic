package repository

import (
	"gorm.io/gorm"
	"haliscicek.com/gin-api/model"
)

type BookRepository interface {
	Insert(b model.Book) model.Book
	Update(b model.Book) model.Book
	Delete(b model.Book)
	GetAll() []model.Book
	FindById(bookId uint64) model.Book
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(dbConnection *gorm.DB) BookRepository {
	return &bookConnection{
		connection: dbConnection,
	}
}

func (db *bookConnection) Insert(b model.Book) model.Book {
	db.connection.Save(&b)
	db.connection.Preload("User").Find(&b)
	return b
}
func (db *bookConnection) Update(b model.Book) model.Book {
	db.connection.Updates(&b)
	db.connection.Preload("User").Find(&b)
	return b
}

func (db *bookConnection) Delete(b model.Book) {
	db.connection.Delete(&b)
}
func (db *bookConnection) FindById(id uint64) model.Book {
	var book model.Book
	db.connection.Preload("User").Find(&book, id)
	return book
}

func (db *bookConnection) GetAll() []model.Book {
	var books []model.Book
	db.connection.Preload("User").Find(&books)
	return books
}
