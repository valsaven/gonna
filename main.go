package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	str "strings"
)

// ParseURL get page url and returns file url
func ParseURL(url string) string {
	var fileURL string

	if str.Contains(url, "instagram") {
		fmt.Println("> Instagram")

		if str.Contains(url, "/p/") {
			// Find the review items
			fmt.Println(">> Photo")
		} else if str.Contains(url, "stories") {
			fmt.Println(">> Stories")

			// Photo
			fmt.Println(">>> Photo")

			// Video
			fmt.Println(">>> Video")
		}
	} else if str.Contains(url, "z0r.de") {
		fmt.Println("> z0r")

		// Get html from the page
		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		source, _ := ioutil.ReadAll(resp.Body)
		html := string(source[:])

		// Filter html to get swf link
		index := str.Index(html, "swfobject.embedSWF")
		dirtyURL := html[index+20 : 1000]
		end := str.Index(dirtyURL, "\"")

		fileURL = dirtyURL[:end]
	} else {
		fmt.Println("Nothing to download...")
	}

	return fileURL
}

// DownloadFile will download a url to a local file.
// It's efficient because it will write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var url, fileURL string

	if len(os.Args) > 1 {
		url = os.Args[1]
		fileURL = ParseURL(url)

		if fileURL != "" {
			err := DownloadFile("\\Downloads\\"+filepath.Base(fileURL), fileURL)

			if err != nil {
				panic(err)
			}
		}
	} else {
		log.Fatal("Required context `url` was not specified.")
		return
	}
}
