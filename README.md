Manga-G 
[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=In%20terminal%20read%20Manga%20fast%20with%20our%20cli&url=https://www.github.com/4cecoder/manga-g&via=github&hashtags=cli,fast,golang,manga,downloader)

 ### About the Project
 - written for your Terminal
 - in the Go programming Language
 - For use (but not limited to) Linux, MacOS, and Windows Systems
 
 <!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
      </ul>
    </li>
     <li><a href="#usage">Usage</a></li>
     <li><a href="#acknowledgments">Acknowledgments</a></li>
     <li><a href="#installation">Installation</a></li>
     <li><a href="#contributing">Contributing</a></li>
     <li><a href="#license">License</a></li>
  </ol>

</details>


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/4cecoder/manga-g">
    <img src="manga-g.png" alt="Logo" width="400" height="200">
  </a>

  <h3 align="center">Best-README-Template</h3>

  <p align="center">
    An awesome CLI program to read Manga via Terminal!
    <br />
    <a href="https://github.com/4cecoder/manga-g/doc/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/othneildrew/Best-README-Template">View Demo</a>
    ·
    <a href="https://github.com/4cecoder/manga-g/issues">Report Bug</a>
    ·
    <a href="https://github.com/4cecoder/manga-g/issues/new/choose">Request Feature</a>
  </p>
</div>






## Prerequisites
 - Go Version 1.18 (Should work with previous go versions)
 - git
 - wget 


## Usage:

run MangaG in terminal with this command

`./MangaG`

Paste url with a similar structure such as: `https://somemangasite.com/1749/185053`


# Shouts out to Similar Manga Projects
<Your manga project github repo LINK here UPON pull request>

- project 1 
- anotha manga thing
- third thing here
 
 
 
 
 
 
 ## Installation
 
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


