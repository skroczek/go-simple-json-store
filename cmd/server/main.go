package main

import (
	"github.com/skroczek/acme-restful/pkg/backend/fs"
	"github.com/skroczek/acme-restful/pkg/router"
	"log"
	"os"
	"path/filepath"

	"github.com/skroczek/acme-restful/pkg/server"
)

func main() {
	var err error
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Printf("root: %s", root)

	//issuer := "http://localhost:8081/realms/master"
	//p, err := rs.NewResourceServerClientCredentials(issuer, "acme-client", "OHfGeYpsgqN8FYI2781yY6V969LL9seL")
	//if err != nil {
	//	log.Fatalf("error creating provider %s", err.Error())
	//}
	//o := oicd.NewOicd(p)
	be := fs.NewFilesystemBackend(root, fs.WithCreateDirs(), fs.WithDeleteEmptyDirs())
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
			//router.WithOICD(o),
			// You can add the JWT auth middleware to protect the server by validating the given JWT against
			// public key. The public key must be in PEM format and be provided in the environment variable
			// ACME_RESTFUL_PUBLIC_KEY
			// You can get a public key from keycloak with the following (fish) command:
			// set -gx ACME_RESTFUL_PUBLIC_KEY (curl http://localhost:8081/realms/dev | jq '.public_key' | tr -d '"')
			//router.WithJWTAuth(),
		),
		server.WithListAll(),
		server.WithGetAll())

	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	s.Run()
}
