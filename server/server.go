package server

import (
	"context"
	"log"
	"net"
	"practice/three/database"
	entities "practice/three/database/entities"
	bookpb "practice/three/proto"

	"gorm.io/gorm"
	"google.golang.org/grpc"
)

type BookServer struct {
	bookpb.UnimplementedBookServiceServer
	db *gorm.DB
}


func bookResponse(db *gorm.DB, id uint) (*bookpb.BookResponse, error) {
	if id == 0 {
		return &bookpb.BookResponse{}, nil
	}
	book := entities.Book{}
	DBRes := db.Joins("Author").First(&book, id)

	if DBRes.Error != nil {
		return &bookpb.BookResponse{}, DBRes.Error 
	}

	return &bookpb.BookResponse{Book: entities.TransferBookToPbEntity(book)}, nil
}

func booksByAuthorResponse(db *gorm.DB, id uint) (*bookpb.BooksResponse, error) {
	books := []entities.Book{}
	DBRes := db.Joins("Author").Where(&entities.Book{AuthorID: id}).Find(&books)

	if DBRes.Error != nil {
		return &bookpb.BooksResponse{}, DBRes.Error 
	}

	pbBooks := make([]*bookpb.Book, 0, len(books))
	for i := 0; i < len(books); i++ {
		pbBookRes, err := bookResponse(db, books[i].ID)
		if err != nil {
			log.Fatalf("Problem with parse book response: %v\n", err)
		}
		pbBooks = append(pbBooks, pbBookRes.Book)
	}

	return &bookpb.BooksResponse{ Books: pbBooks}, nil
}

func (s *BookServer) CreateBook(ctx context.Context, req *bookpb.CreateBookRequest) (*bookpb.BookResponse, error) {
	log.Printf("gRPC CreateBook: %v\n", req)
	pbAuthor := req.GetAuthor()
	author := entities.Author{
		FirstName: pbAuthor.FirstName,
		LastName: pbAuthor.LastName,
	}
	s.db.Where(&author).First(&author)
	if author.ID == 0 {
		s.db.Create(&author)
	}
	
	book := entities.Book{
		Name: req.GetName(),
		Author: author,
	}
	s.db.Create(&book)
	return bookResponse(s.db, book.ID)
}

func (s *BookServer) UpdateBook(ctx context.Context, req *bookpb.UpdateBookRequest) (*bookpb.BookResponse, error) {
	log.Printf("gRPC CreateBook: %v\n", req)
	book := entities.Book{
		ID: uint(req.GetId()),
		Name: req.GetName(),
		AuthorID: uint(req.GetAuthor().GetId()),
	}

	s.db.Save(&book)

	return bookResponse(s.db, book.ID)
}

func (s *BookServer) DeleteBook(ctx context.Context, req *bookpb.DeleteBookRequest) (*bookpb.BookResponse, error) {
	log.Printf("gRPC CreateBook: %v\n", req)

	s.db.Delete(&entities.Book{}, uint(req.GetId()))

	return bookResponse(s.db, uint(0))
}

func (s *BookServer) GetBook(ctx context.Context, req *bookpb.GetBookRequest) (*bookpb.BookResponse, error) {
	log.Printf("gRPC GetBook: %v\n", req)
	return bookResponse(s.db, uint(req.GetId()))
}

func (s *BookServer) GetAuthorBooks(ctx context.Context, req *bookpb.GetAuthorBooksRequest) (*bookpb.BooksResponse, error) {
	log.Printf("gRPC GetAuthorBooks: %v\n", req)
	return booksByAuthorResponse(s.db, uint(req.GetId()))
}

func (s *BookServer) ConnectDatabase() {
	Db := database.OpenConnection() // database.go

	log.Println("Database connected")
	s.db = Db
}

func (s *BookServer) Serve(bind string) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC server error: failure ti bind %v\n", err)
	}

	grpcServer := grpc.NewServer()

	bookpb.RegisterBookServiceServer(grpcServer, s)

	log.Printf("gRPC API server listening on %v\n", bind)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\n", err)
	}
}