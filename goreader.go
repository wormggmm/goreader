package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/logger"
	"github.com/wormggmm/goreader/app"
	"github.com/wormggmm/goreader/epub"
	"github.com/wormggmm/goreader/nav"
)

func main() {
	if len(os.Args) <= 1 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		os.Exit(1)
	}
	debugMode := false
	if len(os.Args) == 3 {
		if os.Args[1] == "-d" || os.Args[1] == "--debug" {
			debugMode = true
		}
	}
	filePath := os.Args[len(os.Args)-1]
	fileDir := filepath.Dir(filePath)
	if debugMode {
	} else {
		fileDir = os.DevNull
	}
	lf := newLogger(fileDir)
	defer lf.Close()
	defer logger.Close()
	rc, err := epub.OpenReader(filePath)
	if err != nil {
		var msg string
		switch err {
		case zip.ErrFormat, zip.ErrAlgorithm, zip.ErrChecksum:
			msg = fmt.Sprintf("cannot unzip contents: %s", err.Error())
		default:
			msg = err.Error()
		}
		fmt.Fprintf(os.Stderr, "Unable to open epub: %s\n", msg)
		os.Exit(1)
	}
	defer rc.Close()
	book := rc.Rootfiles[0]

	a := app.NewApp(book, new(nav.Pager), filePath)
	a.Run()

	if a.Err() != nil {
		fmt.Fprintf(os.Stderr, "Exit with error: %s\n", a.Err().Error())
		os.Exit(1)
	}
	os.Exit(0)

}

func newLogger(logPath string) *os.File {
	if logPath != os.DevNull {
		logPath += "/goreader-debug.log"
	}
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	log := logger.Init("LoggerExample", false, true, lf)
	log.Info("Logger initialized")
	return lf
}
func printUsage() {
	fmt.Fprintln(os.Stderr, "goreader [epub_file]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "-h		print keybindings")
}

func printHelp() {
	fmt.Fprintln(os.Stderr, "Key                  Action")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "q / Esc              Quit")
	fmt.Fprintln(os.Stderr, "k / Up arrow         Scroll up")
	fmt.Fprintln(os.Stderr, "j / Down arrow       Scroll down")
	fmt.Fprintln(os.Stderr, "h / Left arrow       Scroll left")
	fmt.Fprintln(os.Stderr, "l / Right arrow      Scroll right")
	fmt.Fprintln(os.Stderr, "b                    Previous page")
	fmt.Fprintln(os.Stderr, "f                    Next page")
	fmt.Fprintln(os.Stderr, "B                    Previous chapter")
	fmt.Fprintln(os.Stderr, "F                    Next chapter")
	fmt.Fprintln(os.Stderr, "g                    Top of chapter")
	fmt.Fprintln(os.Stderr, "G                    Bottom of chapter")
}
