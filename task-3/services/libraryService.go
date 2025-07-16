package services

import (
	"fmt"
	"libraryManagement/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

func BookFactory() func(title string, author string) models.Book {
	id := 0
	return func(title string, author string) models.Book {
		id += 1
		return models.Book{
			ID:     id,
			Title:  title,
			Author: author,
			Status: "Available",
		}
	}
}

func MemberFactory() func(name string) models.Member {
	id := 0
	return func(name string) models.Member {
		id += 1
		return models.Member{
			ID:   id,
			Name: name,
		}
	}
}

type Library struct {
	Books        map[int]models.Book
	Members      map[int]models.Member
	newBookGen   func(title, author string) models.Book
	newMemberGen func(name string) models.Member
}

func (lib *Library) NewLibrary() *Library {
	return &Library{
		Books:        make(map[int]models.Book),
		Members:      make(map[int]models.Member),
		newBookGen:   BookFactory(),
		newMemberGen: MemberFactory(),
	}
}

func (lib *Library) AddBook(book models.Book) {
	newBook := lib.newBookGen(book.Title, book.Author)
	lib.Books[newBook.ID] = newBook
}

func (lib *Library) RemoveBook(bookID int) {
	delete(lib.Books, bookID)
}

func (lib *Library) BorrowBook(bookID int, memberID int) error {
	book, bookExist := lib.Books[bookID]
	member, memberExist := lib.Members[memberID]

	if !bookExist || !memberExist {
		return fmt.Errorf("error, book or member doesnt exist")
	}
	book.Status = "Borrowed"
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	lib.Members[memberID] = member
	return nil
}

func (lib *Library) ReturnBook(bookID int, memberID int) error {
	book, bookExist := lib.Books[bookID]
	member, memberExist := lib.Members[memberID]

	if !bookExist || !memberExist {
		return fmt.Errorf("error, book or member doesnt exist")
	}

	book.Status = "Available"
	var found bool = false

	idx := 0

	for i, val := range member.BorrowedBooks {
		if val.ID == bookID {
			idx = i
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("error, %v didnt borrow %v with an id of %v", member.Name, book.Title, book.ID)
	}

	member.BorrowedBooks = append(member.BorrowedBooks[:idx], member.BorrowedBooks[idx+1:]...)

	lib.Books[bookID] = book
	lib.Members[memberID] = member

	return nil
}

func (lib *Library) ListAvailableBooks() []models.Book {
	var availableBooks []models.Book
	for _, book := range lib.Books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (lib *Library) ListBorrowedBooks(memberID int) []models.Book {
	var borrowedBooks []models.Book
	member, memberExist := lib.Members[memberID]

	if !memberExist {
		return borrowedBooks
	}

	borrowedBooks = append(borrowedBooks, member.BorrowedBooks...)

	return borrowedBooks
}

func (lib *Library) AddMember(name string) {
	newMember := lib.newMemberGen(name)
	lib.Members[newMember.ID] = newMember
}
