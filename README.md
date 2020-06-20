# goboard
Simple Soundboard Web App written in Go

___

Install:

```bash
$ go get github.com/btnmasher/goboard

```

Run:
```bash
$ goboard

```
___

Goboard runs a local webserver on port 8080.

Scans the `clips/` folder in the working directory where `goboard` was run from and automatically populates the page with buttons (needs refresh when adding/removing clips).

Files must be one of the following: `.mp3 .ogg .wav`