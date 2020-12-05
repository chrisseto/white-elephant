//go:generate go run github.com/kevinburke/go-bindata/go-bindata -pkg server -prefix ../../bin/dist/ -o ./bindata.go -tags !dev ../../bin/dist
//go:generate go run github.com/kevinburke/go-bindata/go-bindata -pkg server -prefix ../../bin/dist/ -o ./bindata_dev.go -tags dev -debug ../../bin/dist
package server

import (
	"github.com/go-chi/chi"
	"net/http"
)

func indexHTML(rw http.ResponseWriter, req *http.Request) {
	data, err := Asset("index.html")
	if err != nil {
		panic(err)
	}

	rw.Write(data)
	rw.WriteHeader(http.StatusOK)
}

func assets(rw http.ResponseWriter, req *http.Request) {
	path := chi.URLParam(req, "path")

	data, err := Asset(path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Write(data)
	rw.WriteHeader(http.StatusOK)
}
