package main

import (
	"fmt"
)

func main() {
	list := getBrowsers()
	b := bestBrowser(list)
	fmt.Println("The best browser for you was chosen. It is", b.name, "!")

	//do some http listening stuffs, launch the chosen browser in localhost:8080 or something
	b.exec("localhost:8080")
}
