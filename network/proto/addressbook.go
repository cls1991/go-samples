package main

import (
	"bufio"
	"fmt"
	"github.com/cls1991/go-samples/network/proto/generated"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fname := os.Args[1]
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found. Creating new file.\n", fname)
		} else {
			log.Fatalln("Error reading file: ", err)
		}
	}

	book := &addressbook.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book: ", err)
	}
	listPeople(os.Stdout, book)

	addr, err := promptForAddress(os.Stdin)
	if err != nil {
		log.Fatalln("Error with address: ", err)
	}
	book.People = append(book.People, addr)

	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("Failed to encode address book: ", err)
	}
	if err := ioutil.WriteFile(fname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book: ", err)
	}
}

func promptForAddress(r io.Reader) (*addressbook.Person, error) {
	p := &addressbook.Person{}

	rd := bufio.NewReader(r)
	fmt.Print("Enter person ID number: ")
	if _, err := fmt.Fscanf(rd, "%d\n", &p.Id); err != nil {
		return p, err
	}

	fmt.Print("Enter person Name: ")
	name, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Name = strings.TrimSpace(name)

	fmt.Print("Enter email address(blank for none): ")
	email, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Email = strings.TrimSpace(email)

	for {
		fmt.Print("Enter a phone number(or leave blank to finish): ")
		phone, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}
		pn := &addressbook.Person_PhoneNumber{
			Number: phone,
		}

		fmt.Print("Is this a mobile, home or work phone? ")
		ptype, err := rd.ReadString('\n')
		if err != nil {
			return p, err
		}
		ptype = strings.TrimSpace(ptype)
		switch ptype {
		case "mobile":
			pn.Type = addressbook.Person_MOBILE
		case "home":
			pn.Type = addressbook.Person_HOME
		case "work":
			pn.Type = addressbook.Person_WORK
		default:
			fmt.Printf("Unknown phone type %q. Using default.\n", ptype)
		}
		p.Phones = append(p.Phones, pn)
	}

	return p, nil
}

func printPerson(w io.Writer, p *addressbook.Person) {
	fmt.Fprintln(w, "Person ID: ", p.Id)
	fmt.Fprintln(w, " Name: ", p.Name)
	if p.Email != "" {
		fmt.Fprintln(w, " Email Address: ", p.Email)
	}

	for _, pn := range p.Phones {
		switch pn.Type {
		case addressbook.Person_MOBILE:
			fmt.Fprint(w, " Mobile phone #: ")
		case addressbook.Person_HOME:
			fmt.Fprint(w, " Home phone #: ")
		case addressbook.Person_WORK:
			fmt.Fprint(w, " Work phone #: ")
		}
		fmt.Fprintln(w, pn.Number)
	}
}

func listPeople(w io.Writer, book *addressbook.AddressBook) {
	for _, p := range book.People {
		printPerson(w, p)
	}
}

/*
output:
	Person ID:  1
	 Name:  cls1991
	 Email Address:  nobody@gmail.com
	 Mobile phone #: 12345678901
	Person ID:  2
	 Name:  faker
	 Email Address:  faker@gmail.com
	 Mobile phone #: 1234567890
*/
