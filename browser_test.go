package kiosk

import (
	"flag"
	"fmt"
	"testing"
)

func TestMain(m *testing.T) {
	t := flag.Bool("t", false, "Run tests for all installed browsers")
	flag.Parse()

	if *t {
		testCompatible("https://google.com")
	} else {
		if flag.Arg(0) == "" {
			m.Error("Please supply an url as argument")
			return
		}

		Run(flag.Arg(0))
	}
}

//Run opens a browser in Kiosk mode
func Run(url string) {
	list := GetInstalled() //get all installed browsers
	b := BestBrowser(list) //get the best one
	b.Exec(url)            //run that browser with an url
}

//example usage: testSome([]string{"firefox", "chromium"}, "https://google.com")
func testSome(testing []string, url string) {
	for _, b := range testing {
		fmt.Println(Browsers[b].Name)
		err := Browsers[b].Exec(url)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func testCompatible(url string) {
	testing := GetInstalled()
	for _, b := range testing {
		fmt.Println(b.Name)
		err := b.Exec(url)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func testAll(url string) {
	for _, b := range Browsers {
		fmt.Println(b.Name)
		err := b.Exec(url)
		if err != nil {
			fmt.Println(err)
		}
	}
}
