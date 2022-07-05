package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
)

type Quote struct {
	Id   string `json:"id"`
	Text string `json:"value"`
}

func main() {
	fmt.Println("Threads: ", runtime.GOMAXPROCS(-1))
	getQuotes(100)
}

func getQuotes(amount int) {
	wg := sync.WaitGroup{}
	m := sync.Map{}

	for i := 0; i < amount; i++ {
		wg.Add(1)
		go func(currentNumber int) {
			quote, _ := getQuote()
			m.Store(currentNumber, quote)
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		v, ok := m.Load(i)
		if ok {
			q := v.(*Quote).Text
			fmt.Println(i, q)
		}
	}
}

func getQuote() (*Quote, error) {
	res, err := http.Get("https://api.chucknorris.io/jokes/random")
	quote := &Quote{}

	if err != nil {
		return quote, err
	}

	json.NewDecoder(res.Body).Decode(quote)

	defer res.Body.Close()

	return quote, nil
}
