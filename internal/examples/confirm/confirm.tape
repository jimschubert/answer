Output confirm.gif
Output confirm.ascii

Set Margin 20
Set MarginFill "#674EFF"
Set BorderRadius 10
Set WindowBar Colorful
Set CursorBlink false

Require go

# Setup
Hide
Type "go build -o confirm-example . && clear"
Enter
Show

# Recording (input style)...
Set TypingSpeed 100ms
Type "./confirm-example"
Enter
Sleep 2
Type "n"
Enter
Sleep 3

Hide
Type "clear"
Enter
Show

# horizontal selectable
Type "./confirm-example -selectable"
Enter
Sleep 1.5
Right 1
Sleep 0.6
Right 1
Sleep 0.6
Enter
Sleep 3

Hide
Type "clear"
Enter
Show

# vertical selectable
Type "./confirm-example -selectable -vertical"
Enter
Sleep 1.5
Up 1
Sleep 0.6
Up 1
Sleep 0.6
Up 1
Sleep 0.6
Enter
Sleep 3

# Cleanup
Hide
Ctrl+C
Type 'rm ./confirm-example'
Enter
