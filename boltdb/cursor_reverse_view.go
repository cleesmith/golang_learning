package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("bolt_cursor_reverse.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("animals"))
		fmt.Println("Bucket stats:")
		fmt.Printf("%+v\n", b.Stats())
		c := b.Cursor()
		// Iterate over items in sorted key order.
		// Starts from the first key/value pair and assigns the k/v variables to the next key/value on each iteration.
		fmt.Println("Ascending:")
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("A %s is %s.\n", k, v)
		}
		// Iterate over items in reverse sorted key order.
		// Starts from the last key/value pair and assigns the k/v variables to the next key/value on each iteration.
		fmt.Println("Descending:")
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			fmt.Printf("A %s is %s.\n", k, v)
		}
		return nil
	})
}
