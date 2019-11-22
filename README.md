# Bash typing game

## What it is ?
**Bash typing game** is a small game that you can run from your terminal. It challenges you during multiple typing games.

## Install
```bash
go get -u github.com/vdlbk/playbtg
```

## Usage 
### Start the game
```bash
playbtg
```

### Help & Options
```bash
playbtg --help

playbtg

Usage:
  playbtg [flags]

Flags:
  -h, --help               help for playbtg
  -n, --number-mode        The words will be replaced by numbers
  -m, --upper-lower-mode   The words will be displayed with a mix of character in uppercase and lowercase
  -u, --upper-mode         The words will be displayed in uppercase
      --version            version for playbtg
```

## Roadmap

* [x]  Create uppercase mode
* [x]  Create mixin uppercase/lowercase mode
* [x]  Create mode with numbers
* [ ]  Compute/Detect errors and display a list of characters on which the user should practice
* [ ]  Create a fire system (based on: rythm, delay between words, time penalty in case of errors...)
* [ ]  Create a mode with sentences
* [ ]  Create a mode with symbols
* [ ]  Create a mode with line codes
* [ ]  Compute more stats (standard deviation, percentile?)
* [ ]  Compute stats on word length
* [ ]  Create a mode with a maximum of word (defined length) in x seconds
* [ ]  Create a mode where the game change the letter during typing. You would have to be careful before pressing Return
* [ ]  Save some records
* [ ]  Add feature to imports his own words
