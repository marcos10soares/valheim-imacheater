# valheim-imacheater
Valheim cheating tool

## Introduction
This tool aims to allow any player to give items to himself, or change any inventory item type or quantity.

**This can be used for multiplayer public servers, no admin privileges needed.**

### Why?
**Life is too short for the grind, Valheim is so much more than grind.**

> With great power comes great responsibility, don't ruin you friends games. :)


## Releases
1.0.0 - Features: Can modify item quantities

## Warnings
> Be careful when using the tool, it's not fully tested and has bugs, **DO NOT modify any "item listing" that does not have a name**, this is a known bug where the tool is picking up items that are not items and will corrupt your character if changed. (**Anyway, your character is always backed up to the `/bckp` folder where you run the cheat tools exe.**)
> Use it at your own risk, try it first with a new character you don't care for.

## Usage

**Video:** [link](http://recordit.co/YoYBGWod7G)

## To Do
1. [x] Automatically get character data
1. [x] Automatically backup character data
1. [x] Successfully modify item quantities
1. [ ] Successfully add new items to inventory  (in progress - already achieved it but not programmatically)
1. [ ] Successfully change current active power (in progress)
1. [ ] Successfully change skill levels
1. [x] Create GUI Interface
1. [ ] Create file with all items names and max limits for each
1. [ ] Add verification for item max limits
1. [ ] Add option to quick fix the game to optimize performance [valheim-increase-fps-performance](https://www.pcgamer.com/valheim-increase-fps-performance/)

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


