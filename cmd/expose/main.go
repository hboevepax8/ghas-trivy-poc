package main

import (
	"fmt"
)

const (
	API_KEY      = "1234567890abcdef"
	DATABASE_URL = "bruce:iugwefgoqwufgoqdiqofghqu@localhost/data"
)

func main() {
	fmt.Printf("API Key: %s\n\n", API_KEY)
	fmt.Printf("Database URL: %s\n\n", DATABASE_URL)
}
