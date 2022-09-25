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
	// You can use this to encrypt the data as rest. But you have to set the passphrase in the environment
	// variable ACME_RESTFUL_PASSPHRASE
	//be = backend.NewEncrypted(be)
	server := pkg.NewServer(be)

	router := pkg.DefaultRouter(server)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}
