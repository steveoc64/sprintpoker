// +build !wasm

package main

import "github.com/vugu/vugu/devutil"

// <img style="position: absolute; top: 50%; left: 50%;" src="https://cdnjs.cloudflare.com/ajax/libs/galleriffic/2.0.1/css/loader.gif">

func newIndexHandler() devutil.StaticContent {
	return devutil.StaticContent(`<!doctype html>
	<html>
	<head>
	<title>Planning Poker</title>
	<meta charset="utf-8"/>
	<!-- styles -->
	<link rel="stylesheet" href="https://bootswatch.com/4/darkly/bootstrap.min.css">
	</head>
	<body>
	<div id="vugu_mount_point">
	<img style="position: absolute; top: 10%; left: 40%;" src="blue-loading-gif-transparent-8.gif">
	</div>
	<script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script> <!-- MS Edge polyfill -->
	<script src="/wasm_exec.js"></script>
	<!-- scripts -->
	<script>
	var wasmSupported = (typeof WebAssembly === "object");
	if (wasmSupported) {
		if (!WebAssembly.instantiateStreaming) { // polyfill
			WebAssembly.instantiateStreaming = async (resp, importObject) => {
				const source = await (await resp).arrayBuffer();
				return await WebAssembly.instantiate(source, importObject);
			};
		}
		var mainWasmReq = fetch("/main.wasm").then(function(res) {
			if (res.ok) {
				const go = new Go();
				WebAssembly.instantiateStreaming(res, go.importObject).then((result) => {
					go.run(result.instance);
				});		
			} else {
				res.text().then(function(txt) {
					var el = document.getElementById("vugu_mount_point");
					el.style = 'font-family: monospace; background: black; color: red; padding: 10px';
					el.innerText = txt;
				})
			}
		})
	} else {
		document.getElementById("vugu_mount_point").innerHTML = 'This application requires WebAssembly support.  Please upgrade your browser.';
	}
	</script>
	</body>
	</html>
	`)
}
