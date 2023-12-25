package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

type pathEntry struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	pathRedirector := func(w http.ResponseWriter, r *http.Request) {
		pathToLookup := r.RequestURI
		if pathToLookup == "/favicon.ico" {
			w.WriteHeader(200)
			return
		}
		if url, ok := pathsToUrls[pathToLookup]; !ok {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}
	}

	return pathRedirector
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func loadYaml(rawYml []byte) ([]pathEntry, error) {
	pathEntries := []pathEntry{}
	return pathEntries, yaml.Unmarshal(rawYml, &pathEntries)
}

func entriesToMap(pathEntries []pathEntry) (pathMap map[string]string) {
	pathMap = make(map[string]string)

	for _, pathEntry := range pathEntries {
		pathMap[pathEntry.Path] = pathEntry.Url
	}

	return
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathEntries, err := loadYaml(yml)
	if err != nil {
		return nil, err
	}

	return MapHandler(entriesToMap(pathEntries), fallback), nil
}
