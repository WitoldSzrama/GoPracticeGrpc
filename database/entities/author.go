package entities

import (
	"github.com/bxcodec/faker/v3"
	bookpb "practice/three/proto"
)

type Author struct {
	ID uint
	FirstName string
	LastName string 
}

func (a Author) CreateFakeData(amount uint) []Author {
	authors := []Author{}
	for i := uint(0); i < amount; i++ {
		authors = append(authors, Author{
			FirstName: faker.FirstName(),
			LastName: faker.LastName(),
		})
	}

	DB.Create(authors)
	return authors
}

func TransferAuthorToPbEntity(e Author) bookpb.Author {
	return bookpb.Author{
		Id: int64(e.ID),
		FirstName: e.FirstName,
		LastName: e.LastName,
	}
}

func  TransferAuthorFromPbEntity(pb *bookpb.Author) Author {
	return Author{
		ID: uint(pb.GetId()),
		FirstName: pb.GetFirstName(),
		LastName: pb.GetLastName(),
	}
}

func NewAuthor() *Author {
	return &Author{}
}