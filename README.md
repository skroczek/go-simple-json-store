# acme-restful

This project can be used to quickly provide a restful-api. It does not support authentication or authorisation, so it should not be used in a production environment.
Currently, only the file system is supported as a backend. This means that all data are stored as JSON files on the hard disk.

## Usage

The first argument must be an existing and writable folder. This is where the JSON files are stored.
```
go run cmd/server/main.go _tmp/
```

## Magic URLs

Two "magic" URLs are offered:
* __all.json
* __list.json

__all.json returns all data in a folder combined as a list.  
__list.json returns all file names as a list of the folder.

