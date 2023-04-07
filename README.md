# GoSimpleJSONStore

This project can be used to quickly provide a restful-api. It does not support authentication or authorisation, so it should not be used in a production environment.
Currently, only the file system is supported as a backend. This means that all data are stored as JSON files on the hard disk.

## Usage

The first argument must be an existing and writable folder. This is where the JSON files are stored.
```
go run cmd/server/main.go _tmp/
```

## Magic URLs

You can configure three optional "magic" URLs:

* __all.json
* __list.json
* __dir.json

### __all.json

This URL returns all data as a JSON object. The keys are the file names without the extension. The values are the content of the files.

### __list.json

This magic URL "__list.json" returns a JSON array list of all file names. To remove file extensions from the
list of file names, include the **withoutExtension** parameter in the URL. The value of the parameter is not evaluated,
only its presence is checked. If the parameter is present, the file extensions are removed from the list and only the
base names of the files are returned. Here are examples of how to include the parameter.

### __dir.json

This magic URL "__dir.json" returns a JSON array list of all directories of the current directory. The directories are
returned as relative paths to the current directory.

### Usage

```golang
package main

import (
	"github.com/skroczek/go-simple-json-store/backend"
	"github.com/skroczek/go-simple-json-store/backend/fs"
	"github.com/skroczek/go-simple-json-store/router"
	"github.com/skroczek/go-simple-json-store/server"
)

func main() {
	be := fs.NewMemory()
	s := server.NewServer(
		server.WithBackend(be),
		// You can additional add the encrypted backend to encrypt the data as rest. But you have to set the passphrase
		// in the environment variable GO_SIMPLE_JSON_STORE_PASSPHRASE
		//server.WithBackend(backend.NewEncrypted(be)),
		server.WithRouterOptions(
			router.WithDefaultCors(true),
		),
		server.WithListAll(),
		server.WithGetAll(),
		server.WithListDir(),
	)
	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	s.Run()
}
```

## Example

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe"}' http://localhost:8080/users/1.json
{"name":"John Doe"}
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"Jane Doe"}' http://localhost:8080/users/2.json
{"name":"Jane Doe"}
$ curl -X POST -H "Content-Type: application/json" -d '{"name":"John Doe","age":42}' http://localhost:8080/users/1.json
{"name":"John Doe","age":42}
$ curl http://localhost:8080/users/1.json
{"name":"John Doe","age":42}
$ curl http://localhost:8080/users/2.json
{"name":"Jane Doe"}
$ curl http://localhost:8080/users/__all.json
[{"name":"John Doe","age":42},{"name":"Jane Doe"}]
$ curl http://localhost:8080/users/__list.json
["1.json","2.json"]
```

## Backends

### File System
Currently, only the file system is supported as a backend. This means that all data are stored as JSON files on the hard disk.

## Encryption at rest

It is possible to save the data encrypted through the Encrypted Backend. The encrypted backend acts as a proxy before 
the actual backend. The passphrase in the environment variable GO_SIMPLE_JSON_STORE_PASSPHRASE is used for encryption.

### Usage

```golang
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/skroczek/go-simple-json-store/backend"
	"github.com/skroczek/go-simple-json-store/backend/fs"
	"github.com/skroczek/go-simple-json-store/router"
	"github.com/skroczek/go-simple-json-store/server"
)

func main() {
	var err error
	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	log.Printf("root: %s", root)

	be := fs.NewFilesystemBackend(root, fs.WithCreateDirs(), fs.WithDeleteEmptyDirs())
	// the encrypted backend acts as a proxy before the actual backend
	// the passphrase is read from the environment variable GO_SIMPLE_JSON_STORE_PASSPHRASE
    // if the passphrase is not set, the backend will panic
	encryptedBackend := backend.NewEncrypted(be)
	s := server.NewServer(
		server.WithBackend(encryptedBackend),
		// You can additional add the encrypted backend to encrypt the data as rest. But you have to set the passphrase
		// in the environment variable GO_SIMPLE_JSON_STORE_PASSPHRASE
		//server.WithBackend(backend.NewEncrypted(be)),
		server.WithRouterOptions(
			router.WithDefaultCors(true),
		),
		server.WithListAll(),
		server.WithGetAll(),
	)
	// By default, it serves on :8080 unless a
	// PORT environment variable was defined.
	s.Run()
}
```

# License
You can find the license in the LICENSE file.