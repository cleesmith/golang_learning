package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("bolt_cursor1.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start a read-write transaction.
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("animals"))

		b := tx.Bucket([]byte("animals"))
		b.Put([]byte("dog"), []byte("fun"))
		b.Put([]byte("cat"), []byte("lame"))
		b.Put([]byte("tiger"), []byte("awesome"))

		c := b.Cursor()
		// Iterate over items in sorted key order. This starts from the
		// first key/value pair and assigns the k/v variables to the
		// next key/value on each iteration.
		// The loop finishes at the end of the cursor when a nil key is returned.
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("A %s is %s.\n", k, v)
		}

		return nil
	})
}
