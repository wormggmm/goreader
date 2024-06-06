package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/logger"
	termbox "github.com/nsf/termbox-go"
	hook "github.com/wormggmm/gohook"
	"github.com/wormggmm/goreader/epub"
	"github.com/wormggmm/goreader/nav"
	"github.com/wormggmm/goreader/parse"
)

type Application interface {
	Forward()
	Back()
	NextChapter()
	PrevChapter()

	PageNavigator() nav.PageNavigator
	Exit()
	Run()
	Err() error
}

type Option struct {
	DebugMode  bool
	NoBlank    bool
	GlobalHook bool
}

// app is used to store the current state of the application.
type app struct {
	book     *epub.Rootfile
	pager    nav.PageNavigator
	chapter  int
	bookPath string
	opt      *Option

	eventCh chan termbox.Event
	err     error

	exitSignal chan bool
}

// NewApp creates an App
func NewApp(b *epub.Rootfile, bookpath string, opt *Option) Application {
	logger.Info("Creating new app:")
	logger.Info("Title:", b.Title)

	logger.Info("Contributor:", b.Contributor)
	logger.Info("Coverage:", b.Coverage)
	logger.Info("Creator:", b.Creator)
	logger.Info("Description:", b.Description)
	logger.Info("Identifier:", b.Identifier)
	logger.Info("Language:", b.Language)
	logger.Info("Metadata:", b.Metadata)
	bookpath = filepath.Dir(bookpath)
	p := new(nav.Pager)
	p.NotBlank = opt.NoBlank
	return &app{pager: p,
		book:       b,
		exitSignal: make(chan bool, 1),
		bookPath:   bookpath, opt: opt,
		eventCh: make(chan termbox.Event, 1),
	}
}

// Run opens a book, renders its contents within the pager, and polls for
// terminal events until an error occurs or an exit event is detected.
func (a *app) Run() {
	if a.err = termbox.Init(); a.err != nil {
		return
	}
	defer termbox.Flush()
	defer termbox.Close()

	keymap, chmap := initNavigationKeys(a)
	if a.opt.GlobalHook {
		hookCh := initKeyHook(a)
		defer close(hookCh)
	}
	go func() {
		for {
			ev := termbox.PollEvent()
			if !a.opt.GlobalHook {
				a.eventCh <- ev
			}
		}
	}()
	if a.err = a.openChapter(); a.err != nil {
		return
	}
	a.restore()
MainLoop:
	for {
		if a.err = a.pager.Draw(); a.err != nil {
			return
		}
		select {
		case <-a.exitSignal:
			break MainLoop
		case ev := <-a.eventCh:
			switch ev.Type {
			case termbox.EventKey:
				logger.Info("action ch:", ev.Ch, " key:", ev.Key)
				if action, ok := keymap[ev.Key]; ok {
					action()
				} else if action, ok := chmap[ev.Ch]; ok {
					action()
				}
				a.record()
			}
		}
	}
}

func (a *app) Err() error {
	return a.err
}

func (a *app) PageNavigator() nav.PageNavigator {
	return a.pager
}

func initKeyHook(a *app) chan hook.Event {
	evChan := hook.Start()
	go func() {
		for hookEv := range evChan {
			if hookEv.Kind != hook.KeyDown {
				continue
			}
			ev := termbox.Event{
				Type: termbox.EventKey,
			}
			str := hook.RawcodetoKeychar(hookEv.Rawcode)
			if len(str) > 0 {
				if len(str) == 1 {
					ev.Ch = rune(str[0])
				} else {
					switch str {
					case "up arrow":
						ev.Key = 65517
					case "down arrow":
						ev.Key = 65516
					case "left arrow":
						ev.Key = 65515
					case "right arrow":
						ev.Key = 65514
					}
				}
			}
			logger.Info("hookEv:", hookEv.Rawcode, " str:", str, " xxxxx:", ev.Ch, " keyChar:", hookEv.Keychar)
			a.eventCh <- ev
		}
		logger.Info("hook exit")
	}()
	return evChan
}
func initNavigationKeys(a Application) (map[termbox.Key]func(), map[rune]func()) {
	keymap := map[termbox.Key]func(){
		// Pager
		termbox.KeyArrowDown:  a.PageNavigator().ScrollDown,
		termbox.KeyArrowUp:    a.PageNavigator().ScrollUp,
		termbox.KeyArrowRight: a.PageNavigator().ScrollRight,
		termbox.KeyArrowLeft:  a.PageNavigator().ScrollLeft,

		// Navigation
		termbox.KeyEsc: a.Exit,
	}

	chmap := map[rune]func(){
		// PageNavigator
		'j': a.PageNavigator().ScrollDown,
		'k': a.PageNavigator().ScrollUp,
		'h': a.PageNavigator().ScrollLeft,
		'l': a.PageNavigator().ScrollRight,
		'g': a.PageNavigator().ToTop,
		'G': a.PageNavigator().ToBottom,

		// Navigation
		'q': a.Exit,
		'f': a.Forward,
		'b': a.Back,
		'F': a.NextChapter,
		'B': a.PrevChapter,
	}

	return keymap, chmap
}

// Exit requests app termination.
func (a *app) Exit() {

	a.exitSignal <- true
}

type Mark struct {
	Chapter int `json:"chapter"`
	ScrollY int `json:"scroll_y"`
}

func (a *app) markFilePath() string {
	markFilePath := filepath.Join(a.bookPath, "."+a.book.Title+".mark")
	return markFilePath
}
func (a *app) restore() {
	markFilePath := a.markFilePath()
	markFile, err := os.OpenFile(markFilePath, os.O_RDONLY, 0644)
	if err != nil {
		logger.Error("Failed to open mark file:", err)
		return
	}
	defer markFile.Close()
	b := make([]byte, 256)
	count, err := markFile.Read(b)
	if err != nil {
		logger.Error("Failed to read mark file:", err)
		return
	}
	mark := &Mark{}
	logger.Info("restore: ", string(b))
	err = json.Unmarshal(b[:count], mark)
	if err != nil {
		logger.Error("Failed to unmarshal mark file:", err)
		return
	}
	a.chapter = mark.Chapter
	a.openChapter()
	a.pager.SetScrollY(mark.ScrollY)
}
func (a *app) record() {
	markFilePath := a.markFilePath()
	markFile, err := os.OpenFile(markFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	logger.Infof("record: chapter=%d, page=%d", a.chapter, a.pager.ScrollY())
	if err != nil {
		logger.Error("Failed to open mark file:", err)
		return
	}
	defer markFile.Close()
	mark := Mark{Chapter: a.chapter, ScrollY: a.pager.ScrollY()}
	b, err := json.Marshal(mark)
	if err != nil {
		logger.Error("Failed to open mark file:", err)
		return
	}
	_, err = markFile.Write(b)
	if err != nil {
		logger.Error("Failed to open mark file:", err)
		return
	}
}

// openChapter opens the current chapter and renders it within the pager.
func (a *app) openChapter() error {
	f, err := a.book.Spine.Itemrefs[a.chapter].Open()
	if err != nil {
		return err
	}
	doc, err := parse.ParseText(f, a.book.Manifest.Items)
	if err != nil {
		return err
	}
	a.pager.SetDoc(doc)

	return nil
}

// Forward pages down or opens the next chapter.
func (a *app) Forward() {
	if a.pager.PageDown() || a.chapter >= len(a.book.Spine.Itemrefs)-1 {
		return
	}

	// We reached the bottom.
	if a.NextChapter(); a.err == nil {
		a.pager.ToTop()
	}
}

// Back pages up or opens the previous chapter.
func (a *app) Back() {
	if a.pager.PageUp() || a.chapter <= 0 {
		return
	}

	// We reached the top.
	if a.PrevChapter(); a.err == nil {
		a.pager.ToBottom()
	}
}

// nextChapter opens the next chapter.
func (a *app) NextChapter() {
	if a.chapter >= len(a.book.Spine.Itemrefs)-1 {
		return
	}

	a.chapter++
	if a.err = a.openChapter(); a.err == nil {
		a.pager.ToTop()
	}
}

// prevChapter opens the previous chapter.
func (a *app) PrevChapter() {
	if a.chapter <= 0 {
		return
	}

	a.chapter--
	if a.err = a.openChapter(); a.err == nil {
		a.pager.ToTop()
	}
}
