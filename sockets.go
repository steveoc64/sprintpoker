// +build !wasm

package main

import (
	"net/http"
	"strings"
)

type socketHandler struct {
	data *DataModel
}

func NewSocketHandler(data *DataModel) *socketHandler {
	return &socketHandler{data: data}
}

func (s *socketHandler) RequestMatch(r *http.Request) bool {
	return strings.HasPrefix(r.URL.String(), "/ws")
}
