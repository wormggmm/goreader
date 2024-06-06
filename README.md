# goreader
## **FORK From taylorskalyo/goreader**
### Changed
#### 2024/06/05(#2)
1. Add -nb (not blank line mode)
2. Use the golang flag package instead os.Args

#### 2024/06/05(#1)
1. change Hot key for Next/Previous Chapter
2. Automatically save/restore the last read place of each book
   1. just save chapter and page scroll_Y 
   2. the mark is named .{BookTitle}.mark, in the same path with the book.
3. add -d param for print debug log into file(same path with the book)



Terminal epub reader

[![Go Report Card](https://goreportcard.com/badge/github.com/wormggmm/goreader)](https://goreportcard.com/report/github.com/wormggmm/goreader)

Goreader is a minimal ereader application that runs in the terminal. Images are displayed as ASCII art. Commands are based on less.

## Installation

``` shell
go install github.com/wormggmm/goreader
```

## Usage

``` shell
goreader [-h] [-d] [-g] [-nb] [epub_file]

# help print
goreader -h

# print debug info to log file, same path of the book
goreader -d [epub_file]

# none blank line mode
goreader -nb [epub_file]

# hook hotkey can without focus
goreader -g [epub_file]
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
