# noisemap demo tape
# Record with: vhs demo.tape

Output demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "TokyoNight"
Set Padding 20
Set TypingSpeed 60ms
Set PlaybackSpeed 0.75

# Launch gitcinema on a repo
Type "gitcinema ."
Sleep 500ms
Enter

# Wait for history to load
Sleep 5s

# Step through a few commits
Sleep 1s
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms

# Jump to first
Type "g"
Sleep 800ms

# Start playback
Type " "
Sleep 6s

# Pause
Type " "
Sleep 1s

# Jump to end
Type "G"
Sleep 1s

# Back to beginning and play again briefly
Type "g"
Sleep 500ms
Type " "
Sleep 3s
Type " "
Sleep 1s

# Switch to detail pane
Tab
Sleep 1.5s

# Switch back
Tab
Sleep 1s

# Open search
Type "/"
Sleep 500ms
Type "feat"
Sleep 1.5s
Enter
Sleep 1s

# Filter by author
Type "f"
Sleep 500ms
Enter
Sleep 1s

# Quit
Type "q"
Sleep 500ms
