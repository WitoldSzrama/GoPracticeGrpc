package entities

import (
	"github.com/bxcodec/faker/v3"
	bookpb "practice/three/proto"
)

type Book struct {
	ID uint
	Name string `gorm:"type:varchar(200);uniqueIndex:idx_title_author"`
	Author Author
	AuthorID uint
}

func (b Book) CreateFakeData(amount uint) []Book {
	books := []Book{}
	author := NewAuthor().CreateFakeData(1)[0]
	for i := uint(0); i < amount; i++ {
		books = append(books, Book{
			Name: faker.FirstName(),
			AuthorID: author.ID,
		})
	}

	DB.Create(books)
	return books
}

func (b *Book) CreateFakeDataWithAuthor(amount uint, author Author) []Book {
	books := []Book{}
	for i := uint(0); i < amount; i++ {
		books = append(books, Book{
			Name: faker.FirstName(),
			AuthorID: author.ID,
		})
	}
	DB.Create(books)
	return books
}

func TransferBookToPbEntity(e Book) *bookpb.Book {
	pbAuthor := TransferAuthorToPbEntity(e.Author)
	return &bookpb.Book{
		Id: int64(e.ID),
		Name: e.Name,
		Author: &pbAuthor,
	}
}

func TransferBookFromPbEntity(e Book, pb *bookpb.Book) Book {
	author := TransferAuthorFromPbEntity(pb.GetAuthor())
	return Book{
		ID: uint(pb.Id),
		Author: author,
		AuthorID: author.ID,
	}
}

func NewBook() *Book {
	return &Book{}
}