package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

// redirectTo writes the Location response header to url and
// set the status code to 301 to trigger a redirect.
func redirectTo(w http.ResponseWriter, url string) {
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.Path]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		redirectTo(w, url)
	}
}

// mappingEntry maps a redirect from request containing Path to URL.
type mappingEntry struct {
	Path string
	URL  string
}

// parseYAMLMapping parses raw YAML mapping to a mappingEntry slice.
func parseYAMLMapping(yml []byte) ([]mappingEntry, error) {
	var entries []mappingEntry
	err := yaml.Unmarshal(yml, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// buildMap constructs a map from path to URL given a mappingEntry slice.
func buildMap(entries []mappingEntry) map[string]string {
	m := make(map[string]string)
	for _, entry := range entries {
		m[entry.Path] = entry.URL
	}
	return m
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	entries, err := parseYAMLMapping(yml)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(entries)
	return MapHandler(pathMap, fallback), nil
}

// parseJSONMapping parses raw JSON mapping to a mappingEntry slice.
func parseJSONMapping(jsn []byte) ([]mappingEntry, error) {
	var entries []mappingEntry
	err := json.Unmarshal(jsn, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
// [
//
//	{
//	  "path": "/some-path",
//	  "url": "https://www.some-url.com/demo"
//	}
//
// ]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	entries, err := parseJSONMapping(json)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(entries)
	return MapHandler(pathMap, fallback), nil
}
