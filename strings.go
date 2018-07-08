package kiosk

const firefoxstring = `
/* * Do not remove the @namespace line -- it's required for correct functioning */
@namespace url("http://www.mozilla.org/keymaster/gatekeeper/there.is.only.xul"); /* set default namespace to XUL */

/*
* Hide tab bar, navigation bar and scrollbars
* !important may be added to force override, but not necessary
*/
#TabsToolbar {visibility: collapse;}
#navigator-toolbox {visibility: collapse;}
#content browser {margin-right: -14px; margin-bottom: -14px;}`
const osascript = `
tell application "Safari"
  activate
  --delay 10
end tell

to clickViewMenuItem(target)
  tell application "System Events"
    tell process "Safari"
      tell menu bar 1
        tell menu bar item "View"
          tell menu "View"
            try
              click menu item target
            end try
          end tell
        end tell
      end tell
    end tell
  end tell
end clickViewMenuItem

clickViewMenuItem("Enter Full Screen")
delay 2
clickViewMenuItem("Hide Bookmarks Bar")
clickViewMenuItem("Hide Tab Bar")
clickViewMenuItem("Hide Status Bar")
clickViewMenuItem("Hide Toolbar Bar")

tell application "System Events"
  tell application process "Safari"
    perform action "AXShowMenu" of tool bar 1 of group 2 of window 1
    key code 125
    keystroke return
  end tell
end tell

say "Ready to rock."
`
