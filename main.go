package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/gocolly/colly"
	"net/http"
	"io"
	"flag"
)

func downloadFile(filepath string, url string) (err error) {

	out, err := os.Create(filepath)
	if err != nil  {
	  return err
	}
	defer out.Close()
  
	resp, err := http.Get(url)
	if err != nil {
	  return err
	}
	defer resp.Body.Close()
  
	if resp.StatusCode != http.StatusOK {
	  return fmt.Errorf("bad status: %s", resp.Status)
	}
  
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
	  return err
	}
  
	return nil
  }
  

func main() {
	i := 0

	l := flag.String("l", "", "link to download")
	n := flag.String("n", "", "name (obligatory)")
	flag.Parse()
	flag.Usage()

	if err := os.MkdirAll(*n+"_images", os.ModePerm); err != nil {
		log.Fatalf("Błąd podczas tworzenia katalogu: %v", err)
	}
	if err := os.MkdirAll(*n+"_video", os.ModePerm); err != nil {
		log.Fatalf("Błąd podczas tworzenia katalogu video: %v", err)
	}
	c := colly.NewCollector(
		colly.AllowedDomains("fapello.com", "cdn.fapello.com"),

	)
	c.AllowURLRevisit = true


	c.OnResponse(func(r *colly.Response) {
		filename := r.Ctx.Get("filename")
		if filename == "" {
			return
		}

		err := r.Save(filename)
		if err != nil {
			log.Printf("Błąd podczas zapisywania pliku %s: %v", filename, err)
			return
		} else {
			fmt.Println("[+] Zapisano plik:", filename)
		}
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		alt := e.Attr("alt")
		if link == "" {
			return
		}
		if strings.Contains(alt, *n) {
		url := e.Request.AbsoluteURL(link)
		ext := strings.ToLower(filepath.Ext(link))
		if ext == ".jpg"{
			i++
			x := strconv.Itoa(i)
			filename := filepath.Join(*n + "_images", *n+x+ext)
			fmt.Println("[+] img found:", url)

			ctx := colly.NewContext()
			ctx.Put("filename", filename)

			err := c.Request("GET", url, nil, ctx, nil)
			if err != nil {
				fmt.Printf("[-] Błąd podczas odwiedzania %s: %v\n", url, err)
				return
			}
		}
	}
		
	})
	c.OnHTML("video, source", func(e *colly.HTMLElement) {

		link := e.Attr("src")
		if link == "" {
			return
		}
		url := e.Request.AbsoluteURL(link)
		ext := strings.ToLower(filepath.Ext(link))
		if ext == ".mp4" {
			fmt.Println(url)
			x := strconv.Itoa(i)
			filename := filepath.Join(*n+x+ext)
			i++
			downloadFile(filepath.Base("") + filename, url)
			
		}
		
	})



	for n:=0; n<10000;n++{
		f := strconv.Itoa(n)
		c.Visit(*l + f + "/")
	}

}
