package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const appName = "svc-payment"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("file .env not found")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_DATABASE"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_TIMEZONE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("open database connection:", err.Error())
	}
	inst, _ := db.DB()
	defer inst.Close()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Service  string `json:"service"`
			Database string `json:"database"`
		}{appName, "down"}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := inst.PingContext(ctx); err == nil {
			response.Database = "up"

		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request method=%s source=%s path=/time\n", r.Method, r.RemoteAddr)

		response := struct {
			Service string `json:"service"`
			Time    string `json:"time"`
		}{appName, time.Now().Format(http.TimeFormat)}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	listen := fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT"))
	log.Println("Starting HTTP service in", listen)

	log.Fatal(http.ListenAndServe(listen, nil))
}
