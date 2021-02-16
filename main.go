package gkeFunc

import (
	"net/http"

	"github.com/shakilbd009/go-gke-api/src/app"
)

func Entry(w http.ResponseWriter, r *http.Request) {
	app.StartApp(w, r)
}
