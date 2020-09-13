# Bash typing game

## What it is ?
**Bash typing game** is a small game you can run from your terminal. It challenges you during multiple typing games.

The default mode gives you 60 words you have to enter correctly. More mode and customization are incoming.

## Install / Update
```bash
go get -u github.com/vdlbk/playbtg
```

## Usage 
### Start the game
```bash
playbtg
# or
cd $GOPATH/src/github.com/vdlbk/playbtg && go run playbtg.go
```

### Help & Options
```
playbtg --help

playbtg

Usage:
  playbtg <file> [flags]

Examples:
playbtg my_file.txt -i --output=/tmp/result

Flags:
  -h, --help                help for playbtg
  -i, --infinite-attempts   You have an infinite numbers of attempts for each words (By default, you only have 1 attempt)
  -s, --no-space-mode       Disable space between words
  -n, --number-mode         The words will be replaced by numbers
  -o, --output string       Specify the folder in which it will save the results into a file (default "Console")
  -m, --upper-lower-mode    The words will be displayed with a mix of character in uppercase and lowercase
  -u, --upper-mode          The words will be displayed in uppercase
  -v, --version             version for playbtg
```

### Use custom words
You can use your very own list of words by running the binary with a valid path to any kind of file which contains text.
Here's a short example :
```bash
playbtg /folder/my_file.txt
```

Where `/folder/my_file.txt` may looks like:
```txt
fire magma electricity
water ice steam
earth metal rock
wind tornado
Science B*tch!
```

## Roadmap

* [x]  Create uppercase mode
* [x]  Create mixin uppercase/lowercase mode
* [x]  Create mode with numbers
* [x]  Change the way to catch the user answer, it should be through a key logger or something like this
* [x]  Compute/Detect errors and display a list of characters on which the user should practice
* [ ]  Create a fire system (based on: rythm, delay between words, time penalty in case of errors...)
* [ ]  Create a mode with sentences
* [ ]  Create a mode with symbols
* [ ]  Create a mode with line codes
* [ ]  Compute more stats (standard deviation, percentile?)
* [ ]  Add [asciigraph](github.com/guptarohit/asciigraph) ??
* [ ]  Compute stats on word length
* [ ]  Create a mode with a maximum of word (defined length) in x seconds
* [ ]  Create a chaos mode which sometimes change the letter you typed randomly
* [ ]  Save some records
* [x]  Add a feature to import his own words
* [ ]  Add a feature to get random wikipedia page
* [x]  Add an option to choose the output (console or file)
* [x]  Add an option to disable space key/ return key between words

## Issues
* [ ]  Stats are not very accurate when this user entered an error (*e.g. spam on space key*)