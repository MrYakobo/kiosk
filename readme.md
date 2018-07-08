# Kiosk
*Open a web browser WITH KIOSK MODE ACTIVATED*

Tired of bundling [Electron](https://github.com/electron/electron) and growing your application binary with 50MB? This package provides an API for opening a web browser in kiosk mode (that is, without navbars), emulating that Electron Single Window Feel®. And yes while it doesn't have the functionality of Electron it certainly makes it Simple to build UIs.

The idea of a "statically linked" browser should really reach the devs of Electron. Imagine requesting a "new window" a la xorg style from your program. This package is a middle ground in a way, using web browsers that are already installed on the user's machine to display your (web)UI.

This package targets the following browsers (on their respective OS):
- Firefox
- Chrome/ium
- Internet Explorer
- Safari

Compatibility table below:

| Browser           | OS      | Tested? |
|-------------------|---------|---------|
| Firefox           | Linux   | Yes     |
| Chrome/ium        | Linux   | Yes     |
| Firefox           | Windows | No      |
| Chrome/ium        | Windows | No      |
| Safari            | Windows | No      |
| Internet Explorer | Windows | No      |
| Firefox           | macOS   | No      |
| Chrome/ium        | macOS   | No      |
| Safari            | macOS   | No      |

# Installation

`go get github.com/MrYakobo/kiosk`

# API
```go
/* Browser has a printable name, a way to execute it and a "rank" (if chrome and ie are both installed,
choose chrome instead of ie because of better support for kiosk mode and other features). */

type Browser struct {
    Name string
    Rank int
    Exec func(url string) error
}

//BestBrowser returns the browser with the highest ranking
func BestBrowser(list []Browser) Browser

//GetInstalled returns slice of Browsers that work on the user's machine, or nil if unknown OS is running
func GetInstalled() []Browser

```

# Example

```go
url := "https://google.com"
list := kiosk.GetInstalled() //get all installed browsers
b := kiosk.BestBrowser(list) //get the best one according to rank
err := b.Exec(url)            //run that browser with an url
if err != nil {
    fmt.Println(err)
}
```