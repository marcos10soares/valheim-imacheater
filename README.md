# valheim-imacheater
Valheim cheating tool

## Introduction
This tool aims to allow any player to give items to himself, or change any inventory item type or quantity.

**This can be used on multiplayer servers, no admin privileges needed.**

### Why?
**Life is too short for the grind, Valheim is so much more than grind.**

> With great power comes great responsibility, don't ruin your friend's games. 😃

## Requirements
Requires Chrome installed

## Releases
- [v1.0.0](https://github.com/marcos10soares/valheim-imacheater/releases/tag/1.0.0) - Features: Can modify item quantities
- [v1.0.1](https://github.com/marcos10soares/valheim-imacheater/releases/tag/1.0.1) - Fixed bug where names with spaces would not load items
- [v1.1.0](https://github.com/marcos10soares/valheim-imacheater/releases/tag/1.1.0) - Fixed bug where some items would not load, **added new feature to reset power cooldown** 🔥

## Warnings
> Be careful when using the tool, it's not fully tested and has bugs, **DO NOT modify any "item listing" that does not have a name**, this is a known bug where the tool is picking up items that are not items and will corrupt your character if changed. (**Anyway, your character is always backed up to the `/bckp` folder where you run the cheat tools exe.**)
> Use it at your own risk, try it first with a new character you don't care for.

## Usage

1. Create a folder somewhere, example on `Desktop` (Note that your character backups will be saved to `/bckp` inside this folder automatically every time you change any items quantities)
2. copy `valheim-cheats.exe` to that folder
3. Run `valheim-cheats.exe` 
4. Modify the quantities you would like and hit save (do not exceed the max quantities of the items or your character will become corrupted) - every time you hit save, a copy of your character data is created inside the folder `/bckp`

**!!IMPORTANT!!: ALWAYS run the program with the game closed. Exit game, run program, relaunch the game**

**Video:** [link](http://recordit.co/YoYBGWod7G)

## Known bugs
- Some special characters in name may cause items not to load.

## To Do
1. [x] Automatically get character data
2. [x] Automatically backup character data
3. [x] Successfully modify item quantities
4. [ ] Successfully add new items to inventory  (in progress - already achieved it but not programmatically)
5. [ ] Successfully change current active power (in progress)
6. [ ] Successfully change skill levels
7. [x] Create GUI Interface
8. [ ] Create file with all items names and max limits for each
9. [ ] Add verification for item max limits
10. [ ] Add option to quick fix the game to optimize performance [valheim-increase-fps-performance](https://www.pcgamer.com/valheim-increase-fps-performance/)

## Screenshots

![modified item quantities](https://github.com/marcos10soares/valheim-imacheater/blob/main/readme-img/1.jpg?raw=true)
![valheim cheat tool](https://github.com/marcos10soares/valheim-imacheater/blob/main/readme-img/demo.gif?raw=true)


### How to build

Make sure you have Golang installed on your computer.

```
git clone https://github.com/marcos10soares/valheim-imacheater.git
cd valheim-imacheater
go run gen.go
go build
```

Instead of `go build` you can also use: `go build -ldflags "-s -w" -trimpath -o valheim-cheats`

#### To Build for Windows on MacOS
```
env GOOS=windows go build -ldflags "-H windowsgui" -o valheim-cheats.exe
```


