Manga-G 
[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=In%20terminal%20read%20Manga%20fast%20with%20our%20cli&url=https://www.github.com/4cecoder/manga-g&via=github&hashtags=cli,fast,golang,manga,downloader)
 ### Terminal Manga Program in written in Go

![Manga G Logo](manga-g.png "MangaG")


## Usage:

run MangaG in terminal with this command

`./MangaG`

Paste url with a similar structure such as: `https://somemangasite.com/1749/185053`

### Requirements to use Manga-G
 - Go Version 1.18 (Should work with previous go versions)
 - git
 - wget 
 
### If you have Go but NOT 1.18+ then

- Modify the `go.mod` file to your currently installed go base version number

- For example, I have go `1.18.1` installed but in the go mod I only need to write `1.18`

### Don't have Go programming language? Let's FIX DAT!

Manual Install (linux) commands

### `wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz`

### `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz`

### `export PATH=$PATH:/usr/local/go/bin`

### `go version`

### For other operating systems:
https://go.dev/doc/install



# Testing Manga-G

## Terminal commands for project setup and run:

### This command gets the Manga-G code "FROM UP OFF DA HUB"

## `git clone https://github.com/4cecoder/manga-g`

### This command takes you into the project folder

## `cd manga-g`

### Checking if there are any problems before running 

## `go mod tidy && go mod vendor`

### This command goes into the run folder and tries to run the program

## `cd cmd/core && go run main.go`

## `go build main.go && mv main MangaG`


# Shouts out to Similar Manga Projects
<Your manga project github repo LINK here UPON pull request>

- project 1 
- anotha manga thing
- third thing here
