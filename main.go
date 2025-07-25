package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	//"os"

	//"github.com/go-redis/redis"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

//	Database connection
//var db *gorm.DB

var db *sql.DB
var redisClient *redis.Client

//structure

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func init() {
	// Connect to PostgreSQL
	// dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	connStr := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("db connected")
	// }
	if err != nil {
		return fmt.Println("failed to open database: %w", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create users table if it doesn't exist
	if err = createUsersTable(); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil

	// sqlStatement := `INSERT INTO users (username, email, password)VALUES ($1, $2, $3)`
	// _, err = db.Exec(sqlStatement, "akash", "akash@gmail.com", "random")
	// if err != nil {
	// 	panic(err)
	// }

	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&User{}) // Migrate the User model

	// Connect to Redis
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Correct way with context
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signup", signup)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func signin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my signin!\n")
}
func dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, dashboard!\n")
}

func createUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	return err
}

// hashPassword hashes a plain text password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func signup(w http.ResponseWriter, r *http.Request) {

	var users User

	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	connStr := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := sql.Open("postgres", connStr)

	// // Validate user input (e.g., email format, password strength)
	// if err := validateUser(users); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	users.Password = string(hashedPassword)
	// Save the user to the database

	sqlStatement := `INSERT INTO users (username, email, password)VALUES ($1, $2, $3)`
	db, err = db.Exec(sqlStatement, users.Username, users.Email, users.Password) // Assuming a function to create user in DB
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Optionally generate and return JWT token here

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
