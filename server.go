// +build !wasm

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vugu/vugu/devutil"
	"golang.org/x/net/websocket"
)

func main() {
	l := ":8081"
	log.Printf("Starting Poker Server at %q", l)

	devMode := os.Getenv("POKER_DEV") != ""

	// create the datamodel and APIs
	dataModel := NewDataModel()
	socketHandler := NewSocketHandler(dataModel)

	// build the routes and the special routes
	mux := devutil.NewMux()

	// API and common routes
	mux.Match(dataModel, dataModel)

	// add websocket handler
	mux.Match(socketHandler, websocket.Handler(func(conn *websocket.Conn) {
		println(time.Now().String(), conn.RemoteAddr().String(), "WebSocket")
		data := dataModel.Listen()
		for {
			select {
			case <-data:
				println("->")
				conn.Write([]byte(time.Now().String()))
			}
		}
	}))

	mux.Match(devutil.NoFileExt, newIndexHandler())

	if devMode {
		wc := devutil.NewWasmCompiler().SetDir(".")
		mux.Exact("/main.wasm", devutil.NewMainWasmHandler(wc))
		mux.Exact("/wasm_exec.js", devutil.NewWasmExecJSHandler(wc))
		mux.Default(devutil.NewFileServer().SetDir("."))
	} else {
		mux.Default(devutil.NewFileServer().SetDir("/var/www/poker"))
	}

	// run the poker server
	log.Fatal(http.ListenAndServe(l, mux))
}
