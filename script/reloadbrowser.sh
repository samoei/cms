#!/bin/bash

# AppleScript to refresh Firefox Developer Edition
osascript <<EOF
tell application "Firefox Developer Edition"
    activate
    tell application "System Events"
        keystroke "r" using command down
    end tell
end tell

# Return to the previous application
# tell application "System Events"
#     keystroke tab using command down
# end tell
EOF
