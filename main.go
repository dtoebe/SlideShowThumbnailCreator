package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/disintegration/imaging"
)

var thumbDir string
var thumbWidth int
var thumbHeight int

func parsePath(pth string) {
	file, err := os.Open(pth)
	if err != nil {
		log.Printf("[ERR] Cannot open path/file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Printf("[ERR] Cannot get path/file info: %v\n", err)
		os.Exit(1)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		files := listDirContents(pth)
		for _, file := range files {
			fullPath := path.Join(pth, file.Name())
			parseFile(fullPath)
		}
	case mode.IsRegular():
		parseFile(pth)
	}
}

func parseFile(fullPath string) {
	if checkExt(fullPath) {
		createThumbs(fullPath)
	} else {
		log.Printf("[INF] Not Image: %s\n", fullPath)
	}
}

func checkExt(pth string) bool {
	supportedTypes := []string{".jpg", ".jepg", ".png", ".gif"}
	for _, typ := range supportedTypes {
		if path.Ext(pth) == typ {
			return true
		}
	}
	return false
}

func createThumbs(pth string) {
	img, err := imaging.Open(pth)
	if err != nil {
		log.Printf("[ERR] Unalbe to open image: %v\n", err)
		return
	}

	thumb := imaging.Resize(img, thumbWidth, thumbHeight, imaging.Box)

	splitPath := strings.Split(pth, "/")
	thumbPath := path.Join(thumbDir, splitPath[len(splitPath)-1])
	if err = imaging.Save(thumb, thumbPath); err != nil {
		log.Printf("[ERR] Unable to save thumbnail image: %v\n", err)
		return
	}
	log.Printf("[INF] Created Thumbnail: %s\n", thumbPath)
}

func listDirContents(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("[ERR] Unable to read directory contents: %v\n", err)
		os.Exit(1)
	}

	return files
}

func parseFlags() (inp, out string, w, h int) {
	inputPath := flag.String("i", "", "Path to image directory or single file")
	outPath := flag.String("o", "", "Path to thumbnail output directory (Needs to exist)")
	width := flag.Int("w", 0, "Thumbnail Width. If 0 then will keep aspect ratio based on the height")
	height := flag.Int("t", 0, "Thumbnail Height. If 0 then will keep aspect ratio based on the width")
	flag.Parse()

	return *inputPath, *outPath, *width, *height
}

func checkTumbDir(pth string) {
	if _, err := os.Stat(pth); os.IsNotExist(err) {
		log.Printf("Unable to get thumbnail directory: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	in, d, w, h := parseFlags()
	thumbDir = d
	thumbWidth = w
	thumbHeight = h
	parsePath(in)
}
