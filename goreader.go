package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/logger"
	"github.com/wormggmm/goreader/app"
	"github.com/wormggmm/goreader/epub"
)

var (
	version   = "v0.0.7"
	helpPrint bool
	opt       = &app.Option{}
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "version:", version)
		fmt.Fprintln(os.Stderr, "Usage:")
		printUsage()
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "When Reading:")
		printHelp()
	}
	flag.BoolVar(&opt.DebugMode, "d", false, "debug mode(debug log in same directory of the book)")
	flag.BoolVar(&opt.NoBlank, "nb", false, "not blank line")
	flag.BoolVar(&opt.GlobalHook, "g", false, "hook hotkey global(can without focus)")
}
func main() {
	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(1)
	}
	flag.Parse()
	if helpPrint {
		printHelp()
		os.Exit(1)
	}
	args := flag.Args()
	if len(args) <= 0 {
		fmt.Fprintln(os.Stderr, "No epub file specified")
		os.Exit(1)
	}
	filePath := args[0]
	fileDir := filepath.Dir(filePath)
	if !opt.DebugMode {
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

	a := app.NewApp(book, filePath, opt)
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
	fmt.Fprintln(os.Stderr, "goreader [-h] [-d] [-g] [-nb] [epub_file]")
	fmt.Fprintln(os.Stderr, "")
}

func printHelp() {
	fmt.Fprintln(os.Stderr, "	Key                  Action")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "	q / Esc              Quit")
	fmt.Fprintln(os.Stderr, "	k / Up arrow         Scroll up")
	fmt.Fprintln(os.Stderr, "	j / Down arrow       Scroll down")
	fmt.Fprintln(os.Stderr, "	h / Left arrow       Scroll left")
	fmt.Fprintln(os.Stderr, "	l / Right arrow      Scroll right")
	fmt.Fprintln(os.Stderr, "	b                    Previous page")
	fmt.Fprintln(os.Stderr, "	f                    Next page")
	fmt.Fprintln(os.Stderr, "	B                    Previous chapter")
	fmt.Fprintln(os.Stderr, "	F                    Next chapter")
	fmt.Fprintln(os.Stderr, "	g                    Top of chapter")
	fmt.Fprintln(os.Stderr, "	G                    Bottom of chapter")
	fmt.Fprintln(os.Stderr, "	Ctrl/Cmd + 1,2,3     Turn on/off global hotkey listener")
	fmt.Fprintln(os.Stderr, "	Mouse Wheel          Scroll like j/h")
	fmt.Fprintln(os.Stderr, "	m + key1,key2,key3   Add bookmark named key1,key2,key3")
	fmt.Fprintln(os.Stderr, "	n + key1,key2,key3   Load bookmark named key1,key2,key3")
}
