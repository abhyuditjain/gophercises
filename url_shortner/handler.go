package url_shortner

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/boltdb/bolt"
	"net/http"
	"os"
)

// MapHandler will return http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths to their corresponding URL
// If path is not present in the map, then the fallback
// http.Handler will be called
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if url, ok := pathsToUrls[request.URL.Path]; ok {
			http.Redirect(writer, request, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(writer, request)
	}
}

type pathUrl struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func YAMLHandler(yamlBytes []byte, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	pathsUrls, e := parseYaml(yamlBytes)
	if e != nil {
		return nil, e
	}
	pathsToUrls := buildPathUrlMap(pathsUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsonBytes []byte, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	pathsUrls, e := parseJson(jsonBytes)
	if e != nil {
		return nil, e
	}
	pathsToUrls := buildPathUrlMap(pathsUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func BoltHandler(db *bolt.DB, fallback http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("pathsUrls"))
			if bucket != nil {
				url := bucket.Get([]byte(request.URL.Path))
				if url != nil {
					http.Redirect(writer, request, string(url), http.StatusFound)
				}
			}
			return nil
		})

		if err != nil {
			panic(err)
		}
		fallback.ServeHTTP(writer, request)
	}
}

func SeedBoltDBFromCsv(dbname string, boltSeedFileName string) error {
	urls, err := readCsv(boltSeedFileName)
	if err != nil {
		fmt.Println("Can not read CSV")
		panic(err)
	}
	pathsUrls := parseCsv(urls)
	return seedUrls(dbname, pathsUrls)
}

func seedUrls(dbname string, pathUrls []pathUrl) error {
	db, err := bolt.Open(dbname, 0644, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("pathsUrls"))
		if err != nil {
			return err
		}
		for _, pu := range pathUrls {
			err := b.Put([]byte(pu.Path), []byte(pu.Url))
			if err != nil {
				fmt.Printf("Error adding %v to bucket\n", pu)
			}
		}
		return nil
	})
	return err
}

func readCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	urls, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func buildPathUrlMap(pathsUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, v := range pathsUrls {
		pathsToUrls[v.Path] = v.Url
	}
	return pathsToUrls
}

func parseJson(jsonBytes []byte) ([]pathUrl, error) {
	var pathsUrls []pathUrl
	err := json.Unmarshal(jsonBytes, &pathsUrls)
	if err != nil {
		return nil, err
	}
	return pathsUrls, nil
}

func parseYaml(yamlBytes []byte) ([]pathUrl, error) {
	var pathsUrls []pathUrl
	err := yaml.Unmarshal(yamlBytes, &pathsUrls)
	if err != nil {
		return nil, err
	}
	return pathsUrls, nil
}

func parseCsv(urls [][]string) []pathUrl {
	pathsUrls := make([]pathUrl, len(urls))
	for i, v := range urls {
		pathsUrls[i] = pathUrl{
			Path: v[0],
			Url:  v[1],
		}
	}
	return pathsUrls
}
