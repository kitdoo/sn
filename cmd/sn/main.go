package main

import (
	"os"

	_ "github.com/kitdoo/sn/internal/gominversion"

	"github.com/joho/godotenv"
)

func init() {
	file := ".development.env"
	_ = godotenv.Load(file, ".env")
}

func main() {
	if err := Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
