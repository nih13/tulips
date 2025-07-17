package main

import (
	"context"
	"database/sql"
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
)

//	Database connection
//var db *gorm.DB
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
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("db connected")
	}

	sqlStatement := `INSERT INTO users (username, email, password)VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, "akash", "akash@gmail.com", "random")
	if err != nil {
		panic(err)
	}

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
