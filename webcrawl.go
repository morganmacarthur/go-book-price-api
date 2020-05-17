package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Lang struct {
	Name string
	Year int
	URL  string
}

func count(name, url string, c chan<- string) {
	start := time.Now()
	r, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprintf("%s: %s", name, err)
		return
	}
	n, _ := io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	dt := time.Since(start).Seconds()
	c <- fmt.Sprintf("%s %d [%.2fs]\n", name, n, dt)
}

func do(f func(Lang)) {
	input, err := os.Open("./lang.json")
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(input)
	for {
		var lang Lang
		err := dec.Decode(&lang)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		f(lang)
	}
}

func main() {

	start := time.Now()
	c := make(chan string)
	n := 0
	do(func(lang Lang) {
		n++
		go count(lang.Name, lang.URL, c)
	})

	timeout := time.After(1 * time.Second)
	for i := 0; i < n; i++ {
		select {
		case result := <-c:
			fmt.Print(result)
		case <-timeout:
			fmt.Print("Timed out\n")
			return
		}
	}

	fmt.Printf("%.2fs total\n", time.Since(start).Seconds())
}
