package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/NightMan-1/golang-nw/build"
	"github.com/NightMan-1/golang-nw/pkg"
	"os"
	"path/filepath"
	"runtime"
)

var (
	app         = "myapp"
	name        = "My Application"
	bin         = "myapp"
	binDir      = "."
	cacheDir    = "."
	nwVersion   = "v0.32.0"
	nwOs        = runtime.GOOS
	nwArch      = runtime.GOARCH
	toolbar     = true
	includesDir = ""
	fullscreen  = false
	frame       = true
)

func main() {
	if runtime.GOOS == "windows" {
		app = app + ".exe"
		bin = bin + ".exe"
	}

	flag.StringVar(&app, "app", app, "Web application to be wrapped by node-webkit.")
	flag.StringVar(&name, "name", name, "Application name.")
	flag.StringVar(&bin, "bin", bin, "Destination file for combined application and node-webkit .nw file (will be placed in binDir directory).")
	flag.StringVar(&binDir, "binDir", binDir, "Destination directory for bin and dependencies.")
	flag.StringVar(&cacheDir, "cacheDir", cacheDir, "Directory to cache node-webkit download.")
	flag.StringVar(&nwVersion, "version", nwVersion, "node-webkit version.")
	flag.StringVar(&nwOs, "os", nwOs, "Target os [linux|windows|darwin].")
	flag.StringVar(&nwArch, "arch", nwArch, "Target arch [386|amd64].")
	flag.BoolVar(&toolbar, "toolbar", toolbar, "Enable toolbar.")
	flag.StringVar(&includesDir, "includesDir", includesDir, "Directory containing additional files to bundle with the .nw file")
	flag.BoolVar(&fullscreen, "fullscreen", fullscreen, "Enable fullscreen mode.")
	flag.BoolVar(&frame, "frame", frame, "Set to false to make window frameless.")
	flag.Parse()

	p := pkg.New(nwVersion, nwOs, nwArch)

	if filepath.Base(bin) != bin {
		panic(fmt.Errorf("bin %q includes a path", bin))
	}

	nw := filepath.Join(binDir, bin+".nw")
	fmt.Printf("Building:\t %s\n", nw)
	if err := nwBuild(nw); err != nil {
		panic(err)
	}

	fmt.Printf("Downloading:\t %s\n", p.Url)
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		panic(err)
	}
	nodeWebkitPath, err := nwDownload(p)
	if err != nil {
		panic(err)
	}

	out := filepath.Join(binDir, bin)
	fmt.Printf("Packaging:\t %s\n", out)
	// Ensure bin directory exists
	if err := os.MkdirAll(binDir, 0755); err != nil {
		panic(err)
	}

	if err := nwPkg(p, nodeWebkitPath, nw); err != nil {
		panic(err)
	}
}

func nwBuild(nw string) error {
	w, err := os.Create(nw)
	if err != nil {
		return err
	}
	defer w.Close()

	zw := zip.NewWriter(w)
	defer zw.Close()

	r, err := os.Open(app)
	if err != nil {
		return err
	}
	defer r.Close()

	bin := filepath.Base(app)
	p := build.Package{Name: name, Bin: bin, Window: build.Window{Title: name, Toolbar: toolbar, Fullscreen: fullscreen, Frame: frame}}

	if err := p.CreateNW(zw, build.DefaultTemplates, r, includesDir); err != nil {
		return err
	}

	return nil
}

func nwDownload(p pkg.Pkg) (string, error) {
	return p.Download(cacheDir)
}

func nwPkg(p pkg.Pkg, nodeWebkitPath string, nw string) error {
	r, err := os.Open(nw)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := p.Package(nodeWebkitPath, r, bin, binDir); err != nil {
		return err
	}

	return nil
}
