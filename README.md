# goreader
## **FORK From wormggmm/goreader**
### Changed
1. Hot key for Next/Previous Chapter
2. Automatically save/restore the last read place of each book
   1. just save chapter and page scroll_Y, if you changed view width....(fixing)
   2. the mark is named .{BookTitle}.mark, in the same path with the book.


Terminal epub reader

[![Go Report Card](https://goreportcard.com/badge/github.com/wormggmm/goreader)](https://goreportcard.com/report/github.com/wormggmm/goreader)

Goreader is a minimal ereader application that runs in the terminal. Images are displayed as ASCII art. Commands are based on less.

## Installation

``` shell
go install github.com/wormggmm/goreader
```

## Usage

``` shell
goreader [epub_file]

# help print
goreader -h

# print debug info to log file, same path of the book
goreader -d [epub_file]
```

### Keybindings

| Key               | Action            |
| ----------------- | ----------------- |
| `q`               | Quit              |
| `k` / Up arrow    | Scroll up         |
| `j` / Down arrow  | Scroll down       |
| `h` / Left arrow  | Scroll left       |
| `l` / Right arrow | Scroll right      |
| `b`               | Previous page     |
| `f`               | Next page         |
| `B`               | Previous chapter  |
| `F`               | Next chapter      |
| `g`               | Top of chapter    |
| `G`               | Bottom of chapter |
