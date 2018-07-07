package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

//a browser has a printable name, a way to execute it and a "rank" (if chrome and ie are both installed, choose chrome instead of ie because of better support for kiosk mode and other features)
//rank 0 is the most preffered rank.
type browser struct {
	name string
	rank int
	exec func(url string) error
}

type browserSlice []browser

func (b browserSlice) Len() int {
	return len(b)
}
func (b browserSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b browserSlice) Less(i, j int) bool {
	return b[i].rank < b[j].rank
}

//all browsers are held in a map, accessible with their 'shorthand' representations
var browsers = map[string]browser{
	"firefox": browser{
		"Firefox",
		1,
		func(url string) error {
			//get root dir path
			root, err := homedir.Expand("~/.mozilla/firefox/")
			if err != nil {
				return err
			}

			//read that dir
			files, er := ioutil.ReadDir(root)
			if er != nil {
				return er
			}

			//try to find a .default folder
			var dir string = ""
			for _, f := range files {
				if strings.Index(f.Name(), ".default") > -1 {
					dir = f.Name()
				}
			}
			if dir == "" {
				return errors.New("ERR: Could not find ~/.mozilla/firefox/ABC123.default")
			}

			root = path.Join(root, "chrome")
			os.MkdirAll(root, os.ModePerm)

			//check if userChrome already exists
			file := path.Join(root, "userChrome.css")
			if _, err := os.Stat(file); os.IsNotExist(err) {
				fmt.Println("INFO: userChrome.css already exists. Appending to file now.")
			}

			//append to userChrome.css
			f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				return err
			}

			defer f.Close()

			if _, err = f.WriteString(firefoxstring); err != nil {
				return err
			}

			return runCmd("firefox", url)
		}},
	"ie": browser{
		"Internet Explorer",
		2,
		func(url string) error {
			return runCmd("iexplore", "-k", url)
		}},
	"chrome": browser{
		"Google Chrome",
		0,
		func(url string) error {
			return runCmd("chrome", "--kiosk", url)
		}},
	"chromium": browser{
		"Chromium",
		0,
		func(url string) error {
			return runCmd("chromium", "--kiosk", url)
		}},
	"safari": browser{
		"Safari",
		2,
		func(url string) error {
			runCmd("safari", url)
			return runCmd("osascript", osascript)
		}}}

func bestBrowser(list browserSlice) browser {
	sort.Sort(list)
	return list[0]
}

//returns slice of browsers that work on the user's machine
//returns nil if unknown OS is runnings
func getBrowsers() []browser {
	switch runtime.GOOS {
	case "darwin":
		return macos()
	case "freebsd", "linux":
		return linux()
	case "windows":
		return windows()
	}
	return nil
}

//TODO: Find more ways to find browsers on linux
func linux() []browser {
	out, err := spawnCmd("xdg-mime query default x-scheme-handler/http") //firefox.desktop
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return []browser{
		browsers[out[:len(out)-8]]}
}

//TODO:
func macos() []browser {
	return []browser{browsers["chrome"]}
}

//TODO:
func windows() []browser {
	return []browser{browsers["chrome"]}
}
