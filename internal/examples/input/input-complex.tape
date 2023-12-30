Output input-complex.gif
Output input-complex.ascii

Set Margin 20
Set MarginFill "#674EFF"
Set BorderRadius 10
Set WindowBar Colorful
Set CursorBlink false

Require go

# Setup
Hide
Type "go build -o input-example . && clear"
Enter
Sleep 2s
Show

# Recording...
Set TypingSpeed 100ms
Type "./input-example -complex"
Set TypingSpeed 650ms
Enter
Sleep 750ms
Type@850ms 'jam'
Sleep 500ms
Set TypingSpeed 200ms
Backspace 3
Sleep 750ms
Type@850ms 'Jameson'
Sleep 500ms
Set TypingSpeed 200ms
Backspace 2
Enter
Sleep 2.5s

# Cleanup
Hide
Ctrl+C
Type 'rm ./input-example'
Enter
