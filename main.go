package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read credentials from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require", user, password, host, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select version()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		if err := rows.Scan(&version); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("version=%s\n", version)

	// calling createCategory
	fmt.Printf("Creating Category in PostgreSQL!\n")

	err = createCategoryTableAndInsertData(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Category table created and data is inserted!\n")

	// Count the rows in the Category table
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM Category;`).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of rows in Category table: %d\n", count)

}

func createCategoryTableAndInsertData(db *sql.DB) error {
	// Create the Category table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Category (
        CategoryID SERIAL PRIMARY KEY,
        Name VARCHAR(50) NOT NULL
    );`)
	if err != nil {
		return fmt.Errorf("error creating Category table: %v", err)
	}

	// Insert 5 elements into the Category table
	categories := []string{"Electronics", "Apparel", "Toys", "Books", "Furniture"}
	for _, name := range categories {
		_, err = db.Exec(`INSERT INTO Category (Name) VALUES ($1);`, name)
		if err != nil {
			return fmt.Errorf("error inserting into Category: %v", err)
		}
	}

	return nil
}
