package main

import (
	"context"
	"log"
	"os"
	bookpb "practice/three/proto"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Author struct {
	ID uint
	FirstName string
	LastName string 
}

type Book struct {
	ID uint
	Name string
	Author Author
}

func logResponse(res *bookpb.BookResponse, err error) {
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if res.Book == nil {
		log.Println("  Book not found")
	} else {
		log.Printf("   response: %v", res.Book)
	}
}

func createBook(client bookpb.BookServiceClient, book Book) *bookpb.Book {
	log.Println("Create book...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	author := bookpb.Author{
		Id: int64(book.Author.ID),
		FirstName: book.Author.FirstName,
		LastName: book.Author.LastName,
	}
	res, err := client.CreateBook(ctx, &bookpb.CreateBookRequest{
		Name: book.Name,
		Author: &author,
	})

	logResponse(res, err)

	return res.Book
}

func getBook(client bookpb.BookServiceClient, id int64) *bookpb.Book {
	log.Println("Get book...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetBook(ctx, &bookpb.GetBookRequest{
		Id: id,
	})

	logResponse(res, err)

	return res.Book
}

func getAuthorBooks(client bookpb.BookServiceClient, id int64) []*bookpb.Book {
	log.Println("Get author books...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetAuthorBooks(ctx, &bookpb.GetAuthorBooksRequest{
		Id: id,
	})

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("response:")
	for i := 0; i< len(res.Books); i++ {
		log.Printf("  item [%v of %v]: %s\n", i+1, len(res.Books), res.Books[i])
	}

	return res.Books
}

func updateBook(client bookpb.BookServiceClient, book Book) *bookpb.Book {
	log.Println("Update book...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	author := bookpb.Author{
		Id: int64(book.Author.ID),
		FirstName: book.Author.FirstName,
		LastName: book.Author.LastName,
	}
	res, err := client.UpdateBook(ctx, &bookpb.UpdateBookRequest{
		Name: book.Name,
		Author: &author,
	})

	logResponse(res, err)

	return res.Book
}

func deleteBook(client bookpb.BookServiceClient, id int64) *bookpb.Book {
	log.Println("Delete book...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeleteBook(ctx, &bookpb.DeleteBookRequest{
		Id: id,
	})

	logResponse(res, err)

	return res.Book
}

// Test all clinets request
func main() {
	err := godotenv.Load()
	
	if err != nil {
		log.Fatalf("Problem with .env file %v\n", err)
	}

	conn, err := grpc.Dial(os.Getenv("Port"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := bookpb.NewBookServiceClient(conn)

	log.Println("---------------")
	book := createBook(client, newTestBook())
	log.Println("---------------")
	clientBook := pbToClientModel(book)
	clientBook.Name = "New updated Name"
	updateBook(client, clientBook)
	log.Println("---------------")
	deleteBook(client, book.GetId())
	log.Println("---------------")
	getBook(client, 1)
	log.Println("---------------")
	getAuthorBooks(client, 1)
}

func newTestBook() Book {
	author := Author{
		FirstName: "TestFirst",
		LastName: "TestLast",
	}
	return Book{
		Author: author,
		Name: "Test Book Name",
	}
}

func pbToClientModel(book *bookpb.Book) Book {
	author := Author{
		FirstName: book.Author.GetFirstName(),
		LastName: book.Author.GetLastName(),
		ID: uint(book.Author.GetId()),
	}

	return Book{
		Name: book.GetName(),
		Author: author,
		ID: uint(book.GetId()),
	}
}