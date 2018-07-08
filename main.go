package kiosk

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

//Browser has a printable name, a way to execute it and a "rank" (if chrome and ie are both installed, choose chrome instead of ie because of better support for kiosk mode and other features). Lower rank is better
type Browser struct {
	Name string
	Rank int
	Exec func(url string) error
}

//BrowserOptions should be passed to .Exec in Browser
// type BrowserOptions struct {
// 	Dimensions []int
// 	FullScreen bool
// }

type browserSlice []Browser

func (b browserSlice) Len() int {
	return len(b)
}
func (b browserSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b browserSlice) Less(i, j int) bool {
	return b[i].Rank < b[j].Rank
}

//Browsers are held in a map, accessible with their 'shorthand' representations
var browsers = map[string]Browser{
	"firefox": Browser{
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

			//try to find a .noui folder
			var dir string = ""
			for _, f := range files {
				if strings.Index(f.Name(), ".noui") > -1 {
					dir = f.Name()
				}
			}
			//make that profile!
			if dir == "" {
				out, er := spawnCmd("firefox", "-CreateProfile", "noui", "-no-remote")
				if er != nil {
					return er
				}
				dir = strings.Replace(out, "Success: created profile 'noui' at '", "", 0)
				dir = dir[0 : len(dir)-1] //cut off last '
			}

			root = path.Join(root, dir, "chrome")
			os.MkdirAll(root, os.ModePerm)

			file := path.Join(root, "userChrome.css")
			var f *os.File = nil

			//check if userChrome already exists
			if _, err := os.Stat(file); !os.IsNotExist(err) {
				str, _ := ioutil.ReadFile(file)
				//if this css code is not already in userChrome.css
				if strings.Index((string)(str), firefoxstring) == -1 {
					//append to userChrome.css
					f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
					if err != nil {
						return err
					}
				}
				// fmt.Println("INFO: userChrome.css already exists. Appending to file now.")

			} else {
				f, err = os.Create(file)
				if err != nil {
					return err
				}
			}

			if f != nil {
				defer f.Close()

				if _, err = f.WriteString(firefoxstring); err != nil {
					return err
				}
			}

			return runCmd("firefox", "-new-instance", "-P", "noui", url)
		}},
	"ie": Browser{
		"Internet Explorer",
		2,
		func(url string) error {
			return runCmd("iexplore", "-k", url)
		}},
	"chrome": Browser{
		"Google Chrome",
		0,
		func(url string) error {
			return runCmd("chrome", strings.Join([]string{"--app", url}, "="))
		}},
	"chromium": Browser{
		"Chromium",
		0,
		func(url string) error {
			return runCmd("chromium", strings.Join([]string{"--app", url}, "="))
		}},
	"safari": Browser{
		"Safari",
		2,
		func(url string) error {
			runCmd("safari", url)
			return runCmd("osascript", osascript)
		}}}

//BestBrowser returns the browser with the highest ranking
func BestBrowser(list []Browser) Browser {
	lst := (browserSlice)(list)
	sort.Sort(lst)
	return lst[0]
}

//GetInstalled returns slice of Browsers that work on the user's machine, or nil if unknown OS is running
func GetInstalled() []Browser {
	switch runtime.GOOS {
	case "freebsd", "linux", "darwin":
		return envLookup(":")
	case "windows":
		return envLookup(";")
	}
	return nil
}

func envLookup(splitter string) []Browser {
	//searches the directories in $PATH
	path := strings.Split(os.Getenv("PATH"), ":")
	browserKeys := reflect.ValueOf(browsers).MapKeys()

	valid := browserSlice{}
	alreadyLogged := map[string]bool{}

	//iterate dirs in $PATH
	for _, d := range path {
		files, _ := ioutil.ReadDir(d)
		//iterate files
		for _, f := range files {
			//iterate Browsers
			for _, b := range browserKeys {
				//if file == Browser
				if f.Name() == b.String() {
					//if not already logged
					if _, ok := alreadyLogged[f.Name()]; !ok {
						alreadyLogged[f.Name()] = true
						valid = append(valid, browsers[b.String()])
					}
				}
			}
		}
	}

	sort.Sort(valid) //sort according to preference
	return valid
}
