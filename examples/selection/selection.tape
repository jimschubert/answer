Output selection.gif
Output selection.ascii

Set Margin 20
Set MarginFill "#674EFF"
Set BorderRadius 10
Set WindowBar Colorful
Set CursorBlink false

Require go

# Setup
Hide
Type "go build -o selection-example . && clear"
Enter
Sleep 2s
Show

# Recording...
Set TypingSpeed 100ms
Type "./selection-example"
Sleep 400ms
Enter
Sleep 1s
Type "?"
Sleep 1.5s
Type "?"
Sleep 1s

Right 2 Sleep 800ms
Left 2 Sleep 800ms

Down 2 Sleep 500ms Space Sleep 500ms
Down 2 Sleep 500ms Space Sleep 500ms

Right 1 Sleep 800ms

Down 3 Sleep 500ms Space Sleep 500ms
Down 1 Sleep 500ms Space Sleep 500ms
Up 1 Sleep 500ms Space Sleep 500ms
Down 1 Sleep 500ms Space Sleep 500ms

Left 1 Sleep 800ms
Right 1 Sleep 800ms
Left 1 Sleep 800ms

Enter
Sleep 2.55s

# Cleanup
Hide
Ctrl+C
Type 'rm ./selection-example'
Enter

