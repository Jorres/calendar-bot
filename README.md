## Calendar bot

This is a fairly simple project.

To launch it:

1. Ensure you have `go` installed, tested on 1.19

2. Install dependencies:

```bash
go mod tidy
```

3. Make sure you have the token for the bot. If you are the right person, you have it :)
   Once you have it, please put it into the `token.txt` file into the same directory as `main.go` file.

4. Run the bot. This will launch the bot and continiously print all the logging info to console:

```bash
go run main.go # 
```

Then, text the @calendarNoteBot (https://t.me/calendarNoteBot) and do:

1. To add the note (make sure you have the ";" symbol in between the date and the actual message):

```
/add 27 April 2023 ; take out the trash
```

2. To view all the notes that you have:

```
/show
```

Enjoy!
