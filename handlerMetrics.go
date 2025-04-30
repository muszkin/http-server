package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	template := "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>"
	_, _ = w.Write([]byte(fmt.Sprintf(template, cfg.fileserverHits.Load())))
}
