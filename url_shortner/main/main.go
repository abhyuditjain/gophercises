package main

import (
	"flag"
	"fmt"
	"github.com/abhyuditjain/gophercices/url_shortner"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"net/http"
)

func main() {
	jsonFilename := flag.String("json", "", "JSON file containing 'path': 'url'")
	yamlFilename := flag.String("yaml", "", "YAML file containing 'path': 'url'")
	boltFilename := flag.String("bolt", "", "BoltDB file. If -seed is also provided, the db created/used will be named this")
	boltSeed := flag.String("seed", "", "BoltDB seed file in CSV format: path,url The db created will be named 'my.db' if -bolt is not provided")

	flag.Parse()

	if *boltSeed != "" {
		if *boltFilename == "" {
			*boltFilename = "my.db"
		}
		err := url_shortner.SeedBoltDBFromCsv(*boltFilename, *boltSeed)
		if err != nil {
			panic(err)
		}
	}

	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort": "https://godoc.org/github.com/gophercises/urlshort",
		"/map":      "https://google.com",
	}

	mapHandler := url_shortner.MapHandler(pathsToUrls, mux)

	if *yamlFilename != "" {
		file, err := ioutil.ReadFile(*yamlFilename)
		if err != nil {
			panic(err)
		}

		yamlHandler, err := url_shortner.YAMLHandler(file, mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		_ = http.ListenAndServe(":8080", yamlHandler)
	} else if *jsonFilename != "" {
		file, err := ioutil.ReadFile(*jsonFilename)
		if err != nil {
			panic(err)
		}

		jsonHandler, err := url_shortner.JSONHandler(file, mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		_ = http.ListenAndServe(":8080", jsonHandler)
	} else if *boltFilename != "" {
		db, err := bolt.Open(*boltFilename, 0600, nil)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		boltHandler := url_shortner.BoltHandler(db, mapHandler)
		fmt.Println("Starting the server on :8080")
		_ = http.ListenAndServe(":8080", boltHandler)
	} else {
		fmt.Println("Starting the server on :8080")
		_ = http.ListenAndServe(":8080", mapHandler)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello, world!")
}
