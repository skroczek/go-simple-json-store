package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/skroczek/acme-restful/pkg/backend"
	"github.com/skroczek/acme-restful/pkg/router"
	"github.com/skroczek/acme-restful/pkg/server"
)

func main() {
	var err error
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Printf("root: %s", root)

	be := backend.NewFilesystemBackend(root)
	s := server.NewServer(
		server.WithBackend(be),
		// You can additional add the encrypted backend to encrypt the data as rest. But you have to set the passphrase
		// in the environment variable ACME_RESTFUL_PASSPHRASE
		//pkg.WithBackend(backend.NewEncrypted(be)),
		server.WithRouterOptions(
			router.WithDefaultCors(true),
			// You can add the basic auth middleware to protect the server with a username and password.
			//router.WithBasicAuth(gin.Accounts{
			//	"admin": "admin",
			//  "user1": "pass1",
			//})
		))

	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	s.Run()
}
