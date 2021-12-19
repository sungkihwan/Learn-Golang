package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	ccsv "github.com/tsak/concurrent-csv-writer"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

func scrape(term string) {
	var baseUrl string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"

	// totalPages := getPages()
	// fmt.Println(strconv.Itoa(totalPages))
	var jobs []extractedJob
	totalPages := 5

	// goroutines
	c := make(chan []extractedJob)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseUrl, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("number of extractedJobs : ", len(jobs))
}

func getPage(page int, baseUrl string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	pageUrl := baseUrl + "&start=" + strconv.Itoa(page*50)
	fmt.Println("start :" + pageUrl)

	// goroutines
	c := make(chan extractedJob)

	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tapItem")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func writeJobs(jobs []extractedJob) {
	var baseIdUrl string = "https://kr.indeed.com/viewjob?jk="
	csv, err := ccsv.NewCsvWriter("sample.csv")
	checkErr(err)

	// Flush pending writes and close file upon exit of main()
	defer csv.Close()

	csv.Write([]string{"Link", "Title", "Location", "Salary", "Summary"})
	done := make(chan bool)
	for _, job := range jobs {
		go func(job extractedJob) {
			csv.Write([]string{baseIdUrl + job.id, job.title, job.location, job.salary, job.summary})
			done <- true
		}(job)
	}

	for i := 0; i < len(jobs); i++ {
		<-done
	}
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".jobTitle>span").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary}
}

func getPages(baseUrl string) int {
	pages := 0
	res, err := http.Get(baseUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection) {
		// fmt.Printf("Review %d: %s\n", i, title)
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
