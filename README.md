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
- [v1.3.0](https://github.com/marcos10soares/valheim-imacheater/releases/tag/v1.3.0) - Added item max quantity verification, updated UI, added new UI button to change equipped items quantity to max quantity, **added new feature to change item levels, like sword lvl255** 🔥
- [v1.2.0](https://github.com/marcos10soares/valheim-imacheater/releases/tag/v1.2.0) - Minor bug fixes, updated UI, added new UI button to change equipped items quantity to max quantity, **added new feature to change equipped power** 
- [v1.1.0](https://github.com/marcos10soares/valheim-imacheater/releases/tag/v1.1.0) - Fixed bug where some items would not load, **added new feature to reset power cooldown** 🔥
- [v1.0.1](https://github.com/marcos10soares/valheim-imacheater/releases/tag/v1.0.1) - Fixed bug where names with spaces would not load items
- [v1.0.0](https://github.com/marcos10soares/valheim-imacheater/releases/tag/1.0.0) - Features: Can modify item quantities

## Warnings
> Be careful when using the tool, there could be bugs that are not yet known. (**Anyway, your character is always backed up to the `/bckp` folder where you run the cheat tools exe.**)
> Use it at your own risk, try it first with a new character you don't care for.

## Usage

1. Create a folder somewhere, example on `Desktop` (Note that your character backups will be saved to `/bckp` inside this folder automatically every time you change any items quantities)
2. copy `valheim-cheats.exe` to that folder
3. Run `valheim-cheats.exe` 
4. Modify the quantities you would like and hit save (do not exceed the max quantities of the items or your character will become corrupted) - every time you hit save, a copy of your character data is created inside the folder `/bckp`

**!!IMPORTANT!!: No need to close the game, just logout of the world you're playing to use the cheat tool**
**Example Usage Video v1.3.0:** [link](https://www.youtube.com/watch?v=LLhJQRQLTpc)
**Example Usage Video v1.2.0:** [link](https://youtu.be/A0wXfo-GB5U)

## Known bugs
- None.

If you find a bug please submit an issue.

## To Do
1. [x] Automatically get character data
2. [x] Automatically backup character data
3. [x] Successfully modify item quantities
4. [x] Successfully reset power cooldown
5. [ ] Successfully add new items to inventory  (in progress - already achieved it but not programmatically)
6. [x] Successfully change item levels
7. [x] Successfully change current active power
8. [ ] Successfully change skill levels 
9. [x] Create GUI Interface
10. [x] Create file with all items names and max limits for each (items_list.json)
11. [x] Add verification for item max limits
12. [ ] Add option to quick fix the game to optimize performance [valheim-increase-fps-performance](https://www.pcgamer.com/valheim-increase-fps-performance/)
13. [ ] Add pre-configured one-click builds (e.g. build for fighting boss, build for resource gathering, build for building) 
14. [ ] Change Life and Stamin to huge values (technically infinite - in progress)

## Screenshots

![modified item quantities](https://github.com/marcos10soares/valheim-imacheater/blob/main/readme-img/1.jpg?raw=true)
![valheim cheat tool](https://github.com/marcos10soares/valheim-imacheater/blob/main/readme-img/demo_v1.3.0.gif?raw=true)


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


