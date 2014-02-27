package main

import (
	"fmt"
	"log"

	augustine "../"
)

func main() {
	db, err := augustine.Open("/tmp/augustine-test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	doc := &augustine.Doc{
		Text: []byte("Hello world!"),
	}
	id, err := db.Put(doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Put doc successfully; id = %d\n", id)

	doc2, err := db.Get(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved document successfully:\n%s\n", string(doc2.Text))

	db.Print()
}
