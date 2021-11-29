package service

import (
	"fmt"
	"github.com/mashingan/smapping"
	"haliscicek.com/gin-api/dto"
	"haliscicek.com/gin-api/model"
	"haliscicek.com/gin-api/repository"
	"log"
)

type BookService interface {
	Insert(b dto.BookCreateDto) model.Book
	Update(b dto.BookUpdateDto) model.Book
	Delete(b model.Book)
	FindById(id uint64) model.Book
	GetAll() []model.Book
	IsAllowedToEdit(userId string, bookId uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepository,
	}
}

func (service *bookService) Insert(b dto.BookCreateDto) model.Book {
	book := model.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	res := service.bookRepository.Insert(book)
	return res
}

func (service *bookService) Update(b dto.BookUpdateDto) model.Book {
	book := model.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	res := service.bookRepository.Update(book)
	return res
}

func (service *bookService) Delete(b model.Book) {
	service.bookRepository.Delete(b)
}

func (service *bookService) FindById(id uint64) model.Book {
	return service.bookRepository.FindById(id)
}

func (service *bookService) GetAll() []model.Book {
	return service.bookRepository.GetAll()
}

func (service *bookService) IsAllowedToEdit(userId string, bookId uint64) bool {
	b := service.bookRepository.FindById(bookId)
	id := fmt.Sprintf("%v", b.UserID)
	return userId == id
}
