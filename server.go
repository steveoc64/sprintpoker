// +build !wasm

package main

import (
	"log"
	"net/http"

	"github.com/vugu/vugu/devutil"
)

func main() {
	l := ":8081"
	log.Printf("Starting Poker Server at %q", l)

	wc := devutil.NewWasmCompiler().SetDir(".")

	// boot up the API
	usersAPI := newUsersAPI()

	// build the routes and the special routes
	mux := devutil.NewMux()
	mux.Exact("/main.wasm", devutil.NewMainWasmHandler(wc))
	mux.Exact("/wasm_exec.js", devutil.NewWasmExecJSHandler(wc))

	// API routes
	mux.Match(usersAPI, usersAPI)

	// default cases
	mux.Match(devutil.NoFileExt, newIndexHandler())
	mux.Default(devutil.NewFileServer().SetDir("."))

	// run the poker server
	log.Fatal(http.ListenAndServe(l, mux))
}
