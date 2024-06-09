# Match-Three Game

[![ruben-match-three-game](https://snapcraft.io/ruben-match-three-game/badge.svg)](https://snapcraft.io/ruben-match-three-game)
[![goreleaser](https://github.com/Ruben9922/match-three-game-cmd/actions/workflows/release.yml/badge.svg)](https://github.com/Ruben9922/match-three-game-cmd/actions/workflows/release.yml)
[![GitHub](https://img.shields.io/github/license/Ruben9922/match-three-game-cmd)](https://github.com/Ruben9922/match-three-game-cmd/blob/master/LICENSE)

A match-three game for the terminal.

[![asciicast](https://asciinema.org/a/662894.svg)](https://asciinema.org/a/662894)

## Features
* Endless and limited moves modes
* Different "symbol sets" - emojis, shapes, letters and numbers
* Show hint (show a possible move)
  * Note: Showing the hint will score no points for that move

## Usage

### Using a binary
Download the latest binary for your OS and architecture from the [releases page](https://github.com/Ruben9922/match-three-game-cmd/releases). Simply extract and run it; no installation needed.

#### Windows
1. Extract the zip archive.
2. Navigate into the folder where you extracted the files.
3. Run `match-three-game.exe`.

#### Linux or macOS
Extract the tar.gz archive using a GUI tool or the command line, e.g.:
```bash
tar -xvzf match-three-game_0.1.0_linux_x86_64.tar.gz --one-top-level
```

Navigate into the folder where you extracted the files, e.g.:
```bash
cd match-three-game_0.1.0_linux_x86_64/
```

Run the program:
```bash
./match-three-game
```

##### Unidentified developer error on macOS
When running the program on macOS for the first time you may get an error saying the app can't be opened as it's from an unidentified developer. You can bypass the error as follows:
1. Locate the `match-three-game` binary in Finder.
2. Control-click the binary, then select Open from the menu.
3. Click Open in the dialog.

This only needs to be done once - in future you can open the app as normal by double-clicking on it.

For more info, please see [this help page](https://support.apple.com/en-gb/guide/mac-help/mh40616/mac) on the Apple website.

### Using Snap (Linux or macOS only)
If using Linux or macOS (with Snap installed), you can install via Snap using either the desktop store or the command line:
```bash
sudo snap install ruben-match-three-game
```

Run the game using the following command:
```bash
ruben-match-three-game
```

## Future Plans
* Add ability to choose number of symbols (fewer symbols would make the game easier)
* Possible other game modes
  * "Clear the board" mode - symbols don't get replenished; game continues until grid is cleared
  * "Bubble" match mode - you can match three or more adjacent symbols in any shape (not necessarily in a row or column as it is currently)
* Timed mode
* High scores (just locally) (?)
* Homebrew and/or Scoop packages (?)
