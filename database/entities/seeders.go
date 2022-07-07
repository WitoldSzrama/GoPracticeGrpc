package entities

import "log"

func SeedBooksForAuthor(amount uint) []Book {
	log.Println("Seeds book for one author")
	author := NewAuthor().CreateFakeData(1)
	book := NewBook()
	return book.CreateFakeDataWithAuthor(amount, author[0])
}