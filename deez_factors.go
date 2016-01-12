package main

import (
  "fmt"
  "log"
  "os"
  "github.com/joho/godotenv"
)

func main() {
  // load environment variables from .env
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  fmt.Println("Get Deez Factors!", os.Getenv("GITHUB_API_KEY"))
}
