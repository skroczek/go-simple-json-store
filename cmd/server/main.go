package main

import (
	"github.com/gin-contrib/cors"
	"github.com/skroczek/acme-restful/pkg"
	"log"
	"os"
	"path/filepath"

	"github.com/skroczek/acme-restful/pkg/backend"
)

func main() {
	var err error
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Printf("root: %s", root)
	be := backend.NewFilesystemBackend(root)
	server := pkg.NewServer(be)

	router := pkg.DefaultRouter(server)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}
