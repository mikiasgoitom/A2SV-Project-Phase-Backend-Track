package controllers

import (
	"bufio"
	"fmt"
	"libraryManagement/models"
	"libraryManagement/services"
	"os"
	"strconv"
	"strings"
)

func Welcome() {
	fmt.Println("************** LIBRARY MANAGEMENT PROGRAM ******************")
}

func Menu(lib *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	option := 0
	for option != 8 {
		fmt.Println("\n1. Add a book")
		fmt.Println("2. Remove a book")
		fmt.Println("3. Borrow a book")
		fmt.Println("4. Return a book")
		fmt.Println("5. List all the available books")
		fmt.Println("6. List all the borrowed books")
		fmt.Println("7. Add a member")
		fmt.Println("8. Quit")
		fmt.Print("choice an operation( 1 - 8 ): ")
		rawOption, _ := reader.ReadString('\n')

		menu, err := strconv.Atoi(strings.TrimSpace(rawOption))
		option = menu

		if err != nil {
			panic("Invalid operations")
		}

		switch option {
		case 1:
			addBook(lib)
		case 2:
			removeBook(lib)
		case 3:
			borrowBook(lib)
		case 4:
			returnBook(lib)
		case 5:
			listAvailableBooks(lib)
		case 6:
			listBorrowedBooks(lib)
		case 7:
			addMember(lib)
		default:
			fmt.Println("invalid input")
		}
	}

}

func addBook(lib *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	var book models.Book

	fmt.Print("\n(Add Book) Enter the book title: ")
	rawTitle, _ := reader.ReadString('\n')
	book.Title = strings.TrimSpace(rawTitle)

	fmt.Print("\nEnter the book author: ")
	rawAuthor, _ := reader.ReadString('\n')
	book.Author = strings.TrimSpace(rawAuthor)

	book.Status = "Available"

	lib.AddBook(book)
}

func removeBook(lib *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n(Remove Book) Enter book id: ")
	rawBId, _ := reader.ReadString('\n')
	bookID, err := strconv.Atoi(strings.TrimSpace(rawBId))

	if err != nil {
		panic("invalid Book id\n")
	}

	found := false

	for _, book := range lib.Books {
		if book.ID == bookID {
			found = true
			break
		}
	}

	if found {
		lib.RemoveBook(bookID)
		fmt.Println("Successfully removed the book.")
	} else {
		fmt.Printf("failed: invalid book id: %v\n", bookID)
	}
}

func borrowBook(lib *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n(Borrow Book) Enter member id: ")
	rawMId, _ := reader.ReadString('\n')
	memberID, err := strconv.Atoi(strings.TrimSpace(rawMId))

	if err != nil {
		panic("invalid member id\n")
	}

	fmt.Print("\nEnter book id: ")
	rawBId, _ := reader.ReadString('\n')
	bookID, err := strconv.Atoi(strings.TrimSpace(rawBId))

	if err != nil {
		panic("invalid Book id\n")
	}

	foundBook := false
	for _, book := range lib.Books {
		if book.ID == bookID {
			foundBook = true
			break
		}
	}

	foundMember := false
	for _, member := range lib.Members {
		if member.ID == memberID {
			foundMember = true
			break
		}
	}

	if foundMember && foundBook {
		book := lib.Books[bookID]
		member := lib.Members[memberID]
		fmt.Printf("%v, with memberid: %v has successfully borrowed the book: %v with book id: %v.", member.Name, memberID, book.Title, bookID)
	} else {
		fmt.Println("failed: invalid book id or member id")
	}
}

func returnBook(lib *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n(Return Book) Enter member id: ")
	rawMId, _ := reader.ReadString('\n')
	memberID, err := strconv.Atoi(strings.TrimSpace(rawMId))

	if err != nil {
		panic("invalid member id\n")
	}

	fmt.Print("\nEnter book id: ")
	rawBId, _ := reader.ReadString('\n')
	bookID, err := strconv.Atoi(strings.TrimSpace(rawBId))

	if err != nil {
		panic("invalid Book id\n")
	}

	foundBook := false
	for _, book := range lib.Books {
		if book.ID == bookID {
			foundBook = true
			break
		}
	}

	foundMember := false
	for _, member := range lib.Members {
		if member.ID == memberID {
			foundMember = true
			break
		}
	}

	if foundMember && foundBook {
		book := lib.Books[bookID]
		member := lib.Members[memberID]
		lib.ReturnBook(bookID, memberID)
		fmt.Printf("%v, with memberid: %v has successfully returned the book: %v with book id: %v.", member.Name, memberID, book.Title, bookID)
	} else {
		fmt.Println("failed: invalid book id or member id")
	}

}

func listAvailableBooks(lib *services.Library) {
	fmt.Println("Available Books: ")
	availableBooks := lib.ListAvailableBooks()
	for i, book := range availableBooks {
		fmt.Printf("%d. %v by %v", i+1, book.Title, book.Author)
	}
	fmt.Println("")
}

func listBorrowedBooks(lib *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n(List of Borrowed Books) Enter member id: ")
	rawMId, _ := reader.ReadString('\n')
	memberID, err := strconv.Atoi(strings.TrimSpace(rawMId))

	if err != nil {
		panic("invalid member id\n")
	}

	borrowedBook := lib.ListBorrowedBooks(memberID)

	for i, book := range borrowedBook {
		fmt.Printf("%d. %v by %v", i+1, book.Title, book.Author)
	}
	fmt.Println("")
}

func addMember(lib *services.Library) {
	reader := bufio.NewReader(os.Stdin)
	var memberName string

	fmt.Print("\n(Add Member) Enter the member name: ")
	rawName, _ := reader.ReadString('\n')
	memberName = strings.TrimSpace(rawName)

	lib.AddMember(memberName)
	fmt.Printf("Successfully added %v as a member.\n", memberName)
}
