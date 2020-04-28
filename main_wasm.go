// +build wasm

package main

import (
	"fmt"
	"strings"

	"flag"
	"syscall/js"

	"github.com/vugu/vugu"
	"github.com/vugu/vugu/domrender"
)

func main() {

	mountPoint := flag.String("mount-point", "#vugu_mount_point", "The query selector for the mount point for the root component, if it is not a full HTML component")
	flag.Parse()

	fmt.Printf("Entering main(), -mount-point=%q\n", *mountPoint)
	defer fmt.Printf("Exiting main()\n")

	buildEnv, err := vugu.NewBuildEnv()
	if err != nil {
		panic(err)
	}

	renderer, err := domrender.New(*mountPoint)
	if err != nil {
		panic(err)
	}
	defer renderer.Release()

	repaintQ := make(chan bool, 1000)
	poker := NewPoker(repaintQ)

	home := js.Global().Get("window").Get("location").Get("href").String()
	fmt.Printf("home is %+v\n", home)
	wsURL := strings.Replace(strings.Replace(home, "http://", "ws://", 1), "#", "", -1) + "ws"

	ws := js.Global().Get("WebSocket")
	fmt.Printf("ws %T %v\n", ws, ws)
	wss := ws.New(wsURL)
	fmt.Printf("wss %T %v\n", wss, wss)

	wss.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("got a msg on the websocket", this.String(), this.Type(), len(args))
		poker.Load()
		return nil
	}))

	// build all the things either due to an event, or
	// something else that manually updates the app state
	go func() {
		for {
			select {
			case <-repaintQ:
				buildResults := buildEnv.RunBuild(poker)
				err = renderer.Render(buildResults)
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	for ok := true; ok; ok = renderer.EventWait() {
		repaintQ <- true
	}

}
