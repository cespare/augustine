package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	augustine "../"
)

const text = `Did you ever fly a kite in bed?
Did you ever walk with ten cats on your head?
Did you ever milk this kind of cow?
Well, we can do it. We know how.
If you never did, you should.
These things are fun and fun is good.`

func main() {
	db, err := augustine.Open("/tmp/augustine-test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, line := range strings.Split(text, "\n") {
		doc := &augustine.Doc{
			Text: []byte(line),
		}
		id, err := db.Put(doc)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully wrote doc; got id %d\n", id)
	}

	db.Print()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		query := scanner.Bytes()
		fmt.Printf("Query: %q\n", string(query))
		results := db.Search(query)
		if len(results) == 0 {
			fmt.Println("No results.")
			continue
		}
		fmt.Println("Results;")
		for i, doc := range db.Search(query) {
			fmt.Printf("(%d) %q\n", i+1, string(doc.Text))
		}
	}
}
