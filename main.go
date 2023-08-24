package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

var wg sync.WaitGroup

func createFolder() string {
	currentTime := time.Now()
	timeFormatted := currentTime.Format("2006_01_02")

	folderName := fmt.Sprintf("%s-%s", timeFormatted, "webpage")

	err := os.Mkdir(folderName, 0755)
	if err != nil {
		fmt.Println("folder already exist", err)
		return folderName
	}

	return folderName
}

func downloadHtml(url string, folderName string, title string, num string) {
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while fetching the web page:", err)
		return
	}
	defer response.Body.Close()

	htmlBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading response body:", err)
		return
	}

	filePath := fmt.Sprintf("%s/%s-%s.html", folderName, title, num) // 폴더 내에 파일 경로 설정
	err = os.WriteFile(filePath, htmlBytes, 0644)
	if err != nil {
		fmt.Println("Error while writing to file:", err)
		return
	}

	fmt.Println("Web page HTML saved to", filePath)
}

func getUrls(rootUrl string) map[string]string {

	response, err := http.Get(rootUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	linkTextMap := make(map[string]string)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}

		text := strings.TrimSpace(s.Text())

		if strings.Contains(link, "article") {
			linkTextMap[link] = text
		}
	})

	return linkTextMap
}

func makeTitle(text string) string {

	substring := text
	if len(text) > 20 {
		substring = text[:20]
	}

	var cleanedSubstring strings.Builder
	for _, char := range substring {
		if unicode.IsLetter(char) {
			cleanedSubstring.WriteRune(char)
		} else if unicode.IsSpace(char) {
			cleanedSubstring.WriteRune('_')
		}
	}

	result := cleanedSubstring.String()

	return result
}

func main() {
	folderName := createFolder()
	if folderName == "" {
		return
	}
	urls := getUrls("https://news.naver.com")

	num := 0
	for link, text := range urls {
		num++
		wg.Add(1)
		go downloadHtml(link, folderName, makeTitle(text), strconv.Itoa(num))
	}

	wg.Wait()

	fmt.Printf("%d개 저장 완료\n", num)
}
