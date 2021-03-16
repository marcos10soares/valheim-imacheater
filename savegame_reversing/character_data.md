# Save Game Reversing information
All the information in this file consist in deductions made from analyzing the binary files.

So far i've seen 3 different save file sizes.

Do not modify the file if save file is around 1kb (just after character creation, play the intro and after this your save should have around 4MB)

**Region starts at: 4194408 (dec) or 0x400068** (for saves that have 4MB)

**Region starts at: 8389395 (dec) or 0x800313** (for saves that have 8MB)

Still don't have a clue of the reason for the file growing from 4MB to 8MB as 90% the file seem to be zeros.

From my analysis the file is divided on the following regions:

* Header
* Huge chunks of mostly zeros, no idea why (90% of the file is in this region) 
* Map markers
* Character Stats
* Inventory
* Discovered items and recipes
* Food Inventory


## Header 
This file size: 8.398.158 bytes
```
06 25 80 00    21 00 00 00    00 00 00 00    1A 00 00 00    C5 00 00 00    7C 08 00 00    02 00 00 00
E6 D9 A0 BC    FF FF FF FF    01 13 8B 42    45 5E FA 11    42 1A 2C 07    45 01 2E D7    51 41 F2 68
9A 42 6A 12    A0 C3 01 BF    8E 64 45 CD    B3 F0 41 D6    0F 08 45 6E    81 5B 45 D9    24 FD 41 C2
69 F7 44 01    6C 02 40 00    04 00 00 00    00 08 00 00    00 00 00 00    00 00 00 00    00 00 00 00
```

`06 25 80 00`: 8398086 dec
8398158 - 8398086 = 72 bytes

Just know that the first uint32 is related to file size.


## Mostly blank chunks

```
00000000 00000000 00000000 00000000 00000000 00000000 
00000000 00000000 00000000 00000000 00000000 00000000 
00000000 00000000 00000000 00000000 00010101 01010101 
01010101 01010101 01000000 00000000 00000000 00000000
00000000 00000000 00000000 00000000 00000000 00000000 
00000000 00000000 00000000 00000000 00000000 00000000
```


## Map Markers Region
Started at 0x400060
```
.  .  .  .     .  .  .  .     .  $  e  n     e  m  y  _     e  i  k  t     h  y  r  .     ¯  .  C  .   
00 00 00 00    17 00 00 00    0E 24 65 6E    65 6D 79 5F    65 69 6B 74    68 79 72 08    AF 84 43 03

.  .  B  ´     k  Â  Ã  .     .  .  .  .     .  `  ý  .     Â  .  .  .     .  .  .  Ò     Ã  .  .  . 
9F 9A 42 B4    6B C2 C3 09    00 00 00 00    00 60 FD 86    C2 00 00 00    00 A0 86 D2    C3 00 00 00

.  .  .  r     i  v  e  r     .  j  !  Â     .  .  .  .     @  .  Â  .     .  .  .  .     .  .  B  l  
00 00 05 72    69 76 65 72    00 6A 21 C2    00 00 00 00    40 8B 20 C2    00 00 00 00    00 0C 42 6C 

a  c  k        f  o  r  e     s  t     .     .  D  .  .     .  .  è  À     .  Ã  .  .     .  .  .  .
61 63 6B 20    66 6F 72 65    73 74 20 96    1A 44 00 00    00 00 E8 C0    8C C3 00 00    00 00 00 0E 

c  o  p  p     e  r  d  e     p  o  s  i     t  X  .  Ë     Ã  .  .  .     .  ø  \  ?     Ä  .  .  .
63 6F 70 70    65 72 20 64    65 70 6F 73    69 74 58 00    CB C3 00 00    00 00 F8 5C    3F C4 00 00

.  .  .  .     f  o  r  s     a  k  e  n        a  l  t     a  r  .  o     E  Ã  .  .     .  .  .  í  ....
00 00 00 0E    66 6F 72 73    61 6B 65 6E    20 61 6C 74    61 72 A0 6F    45 C3 00 00    00 00 00 ED ....
```
Did not dedicate any time to this region, but I'm assuming it has the coordinates and text for each marker. 

`$enemy_eikthyr` must refer to the world variable of where this boss is located.

## Region after player name appears

Example for the Character: "Bjørn"

```
<-- prev |A|  |<---     B     --->||<---                     C                    --->|   |<-- D  -->|
               B  J  ø  ø     R  N
00 00 00 06    42 6A C3 B8    72 6E 6D 1E    F7 D1 00 00    00 00 00 01    1F 20 00 00    18 00 00 00


|<-- E  -->|  |<-- F  -->|    |<-- G  -->|  |<-- H  -->|   |I|| J||<---             K            --->|
                                                                   G  P     _  E  i  k     t  h  y  r
05 03 05 43    05 03 05 43    C3 C3 0C 43    00 63 64 82    45 0A 47 50    5F 45 69 6B    74 68 79 72


|<-- L  -->|  |<--  M  -->|  |<-- N  -->|   |O||<---            P          --->||<---   Q  --->||<-- R 
                                                I  r  o     n  N  a  i     l  s 
00 00 00 00    67 00 00 00    1F 00 00 00    09 49 72 6F    6E 4E 61 69    6C 73 0A 00    00 00 00 00


 -->||<--   X   -->||<--    Y   -->||S||T|  |<---                       U                     --->|| next item -->
C8 42 02 00    00 00 03 00    00 00 00 01    00 00 00 00    00 00 00 6D    1E F7 D1 00    00 00 00 xx
```

File uses Litle Endian format, keep that in mind when reading more than one byte.

- **A.** Player Name length: Example "bjørn", due to special character has an extra byte, 0x06 --> 06 dec
- **B.** Player Name - variable in size (utf8)
- **C.** No clue...
- **D.** No clue... but is always `18 00 00 00`
- **E.** Life max value
- **F.** Related to life, seems to be a timer or percentage?
- **G.** Stamina max value
- **H.** Related to stamina, seems to be a timer or percentage?
- **I.** No clue...
- **J.** Length of the power name, example: for "GP_Eikthyr" --> 0x0A --> 10 dec
- **K.** Name of the Power, example: "GP_Eikthyr", variable length! Note: if `C` is 0, `D` is 0 and `E` does not exist
- **L.** Cooldown timer of the power, microseconds? not sure yet!
- **M.** No clue... but it's always `67 00 00 00`
- **N.** Number of items in inventory, example: 0x1f --> 31 dec
- **O.** Length of the upcoming item name, example: "IronNails" --> 0x09 --> 09 dec
- **P.** Item name, example: "IronNails", variable length! 
- **Q.** Item quantity, uint32
- **R.** Item durability
- **S.** No clue...
- **T.** Item level
- **U.** No clue...

- **X.** X coordinate of item in inventory 
- **Y.** Y coordinate of item in inventory 

**Notes:**
- From the end of `B` to `J` it's always 35 bytes 
- From `L[0]` to `O`, it's 12 bytes

## Region for Recipes, items discovered and achievements

```
<-- prev item region (last item)-->||A||<---  B  --->||C|  |<---  D  
                                                            $  i  t  e     m  _  a  x     e  _  s  t
00 00 00 00    00 00 00 00    00 00 00 AC    00 00 00 0F    24 69 74 65    6D 5F 61 78    65 5F 73 74

    --->||E|  |<---            F              --->|| next item (pattern repeats) -->
o  n  e        $  i  t  e     m  _  c  l     u  b     $     i  t 
6F 6E 65 0A    24 69 74 65    6D 5F 63 6C    75 62 0C 24    69 74 ...
```

- **A.** No clue...
- **B.** Number of entries (items or other entries)
- **C.** Size of entry string
- **D.** Entry string
- **E.** Size of entry string
- **F.** Entry string


## Food Items Region
```
<---       prev recipes and achievements region (last entry)                       ||A|  |<---  B 
$  t  u  t     o  r  i  a     l  _  w  i     s  h  b  o     n  e  _  t     e  x  t        B  e  a  r 
24 74 75 74    6F 72 69 61    6C 5F 77 69    73 68 62 6F    6E 65 5F 74    65 78 74 07    42 65 61 72 


    --->||C|  |<---      D     --->||<---         E    
d  1  0        H  a  i  r     1  4
64 31 30 06    48 61 69 72    31 34 9F 1F    6B 3F 9F 1F    6B 3F 9F 1F    6B 3F BB 60    3C 3F 81 D2 

                               --->||<--   F   -->||G||<---                  H                   --->|
                                                      C     o  o  k  e     d  L  o  x     M  e  a  t
13 3F BD DC    D4 3E 00 00    00 00 03 00    00 00 0D 43    6F 6F 6B 65    64 4C 6F 78    4D 65 61 74


|<-- I  -->|  |<-- J  -->|   |K ||<---             L            --->||<--   M   -->||<--   N   -->||O|
                                 C  o  o     k  e  d  M     e  a  t
70 83 69 41    AC 6D 05 41    0A 43 6F 6F    6B 65 64 4D    65 61 74 00    45 84 41 80    67 46 41 0A

|<---             P            --->||<--   Q   -->||<--   R   -->||<--   S   -->||<--   T   -->||<--   
F  i  s  h     C  o  o  k     e  d
46 69 73 68    43 6F 6F 6B    65 64 67 1A    95 41 E6 A9    25 41 02 00    00 00 0F 00    00 00 0B 00

U -->||<---    Continues until the file ends     ............. 
00 00 D3 C8    22 41 00 00    80 3F 65 00    00 00 D1 40    94 40 00 00    00 00 64 00    00 00 DC D0
```

- **A.** Upcoming String Size
- **B.** Beard selected for the character
- **C.** Upcoming String Size
- **D.** Hair selected for the character
- **E.** No clue...
- **F.** Number of food items equipped
- **G.** String size of first food item
- **H.** Food item name
- **I.** Health modifier (not sure if it is summed, but think so) - modifying this value will affect character stats region (**E**)
- **J.** Stamina modifier (not sure if it is summed, but think so) - modifying this value will affect character stats region (**G**)
- **K.** String size of second food item
- **L.** Food item name
- **M.** same as **I**
- **N.** same as **J**
- **O.** String size of third food item
- **P.** Food item name
- **Q.** same as **I**
- **R.** same as **J**
- **S.** No clue...
- **T.** No clue...
- **U.** No clue...