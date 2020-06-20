package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image/color"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hansrodtang/randomcolor"
)

var workingDir string

type clip struct {
	File  string
	Name  string
	Color string
}

func main() {
	w, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	workingDir = w

	http.HandleFunc("/", serveBoard)

	fs := http.FileServer(http.Dir("clips"))
	http.Handle("/clips/", http.StripPrefix("/clips/", fs))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}

func loadTemplate(file string) (t *template.Template, err error) {
	t, err = template.ParseGlob(file)
	if err != nil {
		log.Print("Error Parsing Template: ", err.Error())
	}
	return
}

func applyTemplate(tmpl *template.Template, clips *[]clip) (buffer *bytes.Buffer, err error) {
	buffer = new(bytes.Buffer)

	err = tmpl.Execute(buffer, clips)
	if err != nil {
		log.Print(err)
	}
	return
}

func serveBoard(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	clips, err := getClips("clips")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}

	tmpl, err := loadTemplate("soundboard.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}

	page, err := applyTemplate(tmpl, clips)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}

	_, err = w.Write(page.Bytes())
	if err != nil {
		log.Print(err)
	}
}

func getClips(directory string) (*[]clip, error) {

	clips := &[]clip{}

	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		path = strings.Replace(path, directory, "", -1)
		if len(path) > 0 {
			if !isSupported(path) {
				return nil
			}

			name := properTitle(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
			file := strings.Replace(filepath.Join(directory, path), "\\", "/", -1)
			clr := randomcolor.New(randomcolor.Random, randomcolor.DARK)
			rgba := color.RGBAModel.Convert(clr).(color.RGBA)
			cstr := fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
			aclip := clip{File: file, Name: name, Color: cstr}
			*clips = append(*clips, aclip)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return clips, nil
}

func isSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".mp3" || ext == ".ogg" || ext == ".wav"

}

func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " a an on the to "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

func logRequest(req *http.Request) {
	now := time.Now()
	log.Printf("%s - %s [%s] \"%s %s %s\" ",
		req.RemoteAddr,
		"",
		now.Format("02/Jan/2006:15:04:05 -0700"),
		req.Method,
		req.URL.RequestURI(),
		req.Proto)
}
