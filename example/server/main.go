package main

import (
	"github.com/skroczek/go-simple-json-store/backend/fs"
	"github.com/skroczek/go-simple-json-store/router"
	"github.com/skroczek/go-simple-json-store/server"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var err error
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Printf("root: %s", root)

	//issuer := "http://localhost:8081/realms/master"
	//p, err := rs.NewResourceServerClientCredentials(issuer, "go-simple-json-store-client", "OHfGeYpsgqN8FYI2781yY6V969LL9seL")
	//if err != nil {
	//	log.Fatalf("error creating provider %s", err.Error())
	//}
	//o := oidc.NewOidc(p)
	be := fs.NewFilesystemBackend(root, fs.WithCreateDirs(), fs.WithDeleteEmptyDirs())
	s := server.NewServer(
		server.WithBackend(be),
		// You can additional add the encrypted backend to encrypt the data as rest. But you have to set the passphrase
		// in the environment variable GO_SIMPLE_JSON_STORE_PASSPHRASE
		//pkg.WithBackend(backend.NewEncrypted(be)),
		server.WithRouterOptions(
			router.WithDefaultCors(true),
			// You can add the basic auth middleware to protect the server with a username and password.
			//router.WithBasicAuth(gin.Accounts{
			//	"admin": "admin",
			//  "user1": "pass1",
			//})
			//router.WithOIDC(o),
			// You can add the JWT auth middleware to protect the server by validating the given JWT against
			// public key. The public key must be in PEM format and be provided in the environment variable
			// GO_SIMPLE_JSON_STORE_PUBLIC_KEY
			// You can get a public key from keycloak with the following (fish) command:
			// set -gx GO_SIMPLE_JSON_STORE_PUBLIC_KEY (curl http://localhost:8081/realms/dev | jq '.public_key' | tr -d '"')
			//router.WithJWTAuth(),
		),
		server.WithListAll(),
		server.WithGetAll(),
		server.WithListDir())

	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	s.Run()
}
