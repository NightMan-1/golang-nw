/*
Call a golang web application from node-webkit to get a native looking application.


Instructions


Go get golang-nw:

    go get github.com/NightMan-1/golang-nw/cmd/golang-nw-pkg

Create an app:

See https://github.com/NightMan-1/golang-nw/blob/master/cmd/example/main.go
	package main

	import (
		"fmt"
		"github.com/NightMan-1/golang-nw"
		"net/http"
	)

	func main() {
		// Setup our handler
		http.HandleFunc("/", hello)

		// Create a link back to node-webkit using the environment variable
		// populated by golang-nw's node-webkit code
		nodeWebkit, err := nw.New()
		if err != nil {
			panic(err)
		}

		// Pick a random localhost port, start listening for http requests using default handler
		// and send a message back to node-webkit to redirect
		if err := nodeWebkit.ListenAndServe(nil); err != nil {
			panic(err)
		}
	}

	func hello(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from golang.")
	}


Build your app:

    go install .\src\github.com\NightMan-1\golang-nw\cmd\example


Wrap it in node-webkit:

    .\bin\golang-nw-pkg.exe -app=.\bin\example.exe -name="My Application" -bin="myapp.exe" -toolbar=false

    Building:        myapp.exe.nw
    Downloading:     http://dl.node-webkit.org/v0.32.0/nwjs-v0.32.0-win-x64.zip
    Packaging:       myapp.exe

You are now good to go:

    .\myapp.exe

You may want to create your own build script so you can control window dimensions etc.
See http://godoc.org/github.com/NightMan-1/golang-nw/build and
https://github.com/NightMan-1/golang-nw/blob/master/cmd/golang-nw-pkg/pkg.go

Command line options:

    $ ./bin/golang-nw-pkg -h
    Usage of ./bin/golang-nw-pkg:
      -app="myapp": Web application to be wrapped by node-webkit.
      -arch="amd64": Target arch [386|amd64].
      -bin="myapp": Destination file for combined application and node-webkit .nw file (will be placed in binDir directory).
      -binDir=".": Destination directory for bin and dependencies.
      -cacheDir=".": Directory to cache node-webkit download.
      -frame=true: Set to false to make window frameless.
      -fullscreen=false: Enable fullscreen mode.
      -includesDir="": Directory containing additional files to bundle with the .nw file
      -name="My Application": Application name.
      -os="linux": Target os [linux|windows|darwin].
      -toolbar=true: Enable toolbar.
      -version="v0.9.2": node-webkit version.

Known issues:

1) libudev.so.0 - On ubuntu >=13.10 (and similar) libudev.so.0 has been removed.
tl;dr:

    $ ./bin/golang-nw-pkg -app=./bin/example -name="My Application" -bin="myapp" -toolbar=false
    $ sed -i -e 's/udev\.so\.0/udev.so.1/g' myapp

Node-webkit has various work arounds described at https://github.com/rogerwang/node-webkit/wiki/The-solution-of-lacking-libudev.so.0

2) Download of node-webkit appears to stall - It's a ~43MB download and can take longer than expected (it could do with some feedback).

*/
package nw
