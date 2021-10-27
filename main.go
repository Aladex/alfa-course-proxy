package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// course link: "https://alfabank.ru/api/v1/scrooge/currencies/alfa-rates?currencyCode.in=EUR&rateType.eq=makeCash&lastActualForDate.eq=true&clientType.eq=standardCC&date.lte=")

var coursesResponse []byte

func alfaCoursesResponse(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(coursesResponse)
}

func dateGenerate() string {
	t := time.Now()
	return t.Format("2006-01-02T03:04:05-07:00")
}

func courseGetTicker(contentResp *[]byte) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Try to get last courses")
				*contentResp = getCourses(os.Getenv("ALFA_LINK"))
			}
		}
	}()
}

func getCourses(url string) []byte {
	client := http.Client{}

	alfaUrl := url + dateGenerate()
	req, err := http.NewRequest("GET", alfaUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Referer":      []string{"https://alfabank.ru/currency/"},
		"User-Agent":   []string{"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:93.0) Gecko/20100101 Firefox/93.0"},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bodyBytes
}

func main() {
	go courseGetTicker(&coursesResponse)
	http.HandleFunc("/courses", alfaCoursesResponse)
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":8090", nil)
}
