# How to implement everything, high level

# Firefox

```bash
~/.mozilla/firefox/ID.default/ mkdir chrome
cp data/userContent.css .
firefox $URL
```
exec
# Chrome

`chrome --kiosk $URL`

# Internet Explorer

`iexplore -k $URL`

# Safari

`osascript`
`safari`

## Checking for browsers
Have `x-default-browser` (npm) as a reference.