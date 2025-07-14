package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
	//"github.com/go-redis/redis"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	//"gorm.io/gorm"
)

//	Database connection
var db *gorm.DB
var redisClient *redis.Client

func init() {
	// Connect to PostgreSQL
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	//db.AutoMigrate(&User{}) // Migrate the User model

	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // No password set for this example
		DB:       0,  // Use default database
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/signin", signin)
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
