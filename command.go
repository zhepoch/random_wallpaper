package main

var GetDesktopCountCommand = `tell application "System Events" to copy count of desktops to stdout`

var ApplyDesktopCommand = `
tell application "System Events"
	tell desktop %v
		set picture to "%v"
	end tell
end tell
`
