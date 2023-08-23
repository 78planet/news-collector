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

var wg sync.WaitGroup // WaitGroup 선언

func createFolder() string {
	// 현재 시간 정보 가져오기
	currentTime := time.Now()
	timeFormatted := currentTime.Format("2006_01_02")

	// 폴더 이름 생성
	folderName := fmt.Sprintf("%s-%s", timeFormatted, "webpage")

	// 폴더 생성
	err := os.Mkdir(folderName, 0755)
	if err != nil {
		fmt.Println("folder already exist", err)
		return folderName
	}

	return folderName
}

func downloadHtml(url string, folderName string, title string, num string) {
	defer wg.Done()
	// HTTP GET 요청 보내기
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while fetching the web page:", err)
		return
	}
	defer response.Body.Close()

	// 응답 본문 읽기
	htmlBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading response body:", err)
		return
	}

	// HTML 내용을 파일로 저장
	filePath := fmt.Sprintf("%s/%s-%s", folderName, title, num) // 폴더 내에 파일 경로 설정
	err = os.WriteFile(filePath, htmlBytes, 0644)
	if err != nil {
		fmt.Println("Error while writing to file:", err)
		return
	}

	fmt.Println("Web page HTML saved to", filePath)
}

func getUrls(rootUrl string) map[string]string {

	//	// HTTP 요청 보내고 응답 받기
	response, err := http.Get(rootUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// 응답 바디를 GoQuery로 읽기
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 링크와 제목 쌍을 저장할 맵
	linkTextMap := make(map[string]string)

	// 링크와 제목 쌍 찾기
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
	// 앞에서 10글자 추출

	// 특수문자 및 뛰어쓰기 처리
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

