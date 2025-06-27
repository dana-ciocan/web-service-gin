package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Album represents the album structure for database operations
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// InitDB initializes the SQLite database
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./albums.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to SQLite database")

	// Create the albums table if it doesn't exist
	createTable()
	
	// Insert sample data if the table is empty
	insertSampleData()
}

// createTable creates the albums table
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS albums (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		artist TEXT NOT NULL,
		price REAL NOT NULL
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

// insertSampleData adds sample albums if the table is empty
func insertSampleData() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM albums").Scan(&count)
	if err != nil {
		log.Fatal("Failed to count albums:", err)
	}

	if count == 0 {
		sampleAlbums := []Album{
			{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
			{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
			{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
		}

		for _, album := range sampleAlbums {
			err := CreateAlbum(album)
			if err != nil {
				log.Printf("Failed to insert sample album %s: %v", album.ID, err)
			}
		}
		log.Println("Sample data inserted")
	}
}

// GetAllAlbums retrieves all albums from the database
func GetAllAlbums() ([]Album, error) {
	rows, err := DB.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

// GetAlbumByID retrieves a specific album by ID
func GetAlbumByID(id string) (*Album, error) {
	var album Album
	err := DB.QueryRow("SELECT id, title, artist, price FROM albums WHERE id = ?", id).
		Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Album not found
		}
		return nil, err
	}

	return &album, nil
}

// CreateAlbum adds a new album to the database
func CreateAlbum(album Album) error {
	query := "INSERT INTO albums (id, title, artist, price) VALUES (?, ?, ?, ?)"
	_, err := DB.Exec(query, album.ID, album.Title, album.Artist, album.Price)
	return err
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}