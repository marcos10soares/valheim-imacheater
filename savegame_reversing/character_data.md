# Save Game Reversing information
All the information in this file consist in deductions made from analysing the binary files.

## Region after player name appears

So far i've seen 3 different save file sizes.

Do not modify file if save file is around 1kb (just after character creation, play the intro and after this your save should have around 4MB)

**Region starts at: 4194408 (dec) or 0x400068** (for saves that have 4MB)

**Region starts at: 8389395 (dec) or 0x800313** (for saves that have 8MB)

Example for the Character: "Bjørn"

```
|<--  A -->|  |<---     B     --->|
               B  J  ø  ø     R  N
00 00 00 06    42 6A C3 B8    72 6E 6D 1E    F7 D1 00 00    00 00 00 01    1F 20 00 00    18 00 00 00


                                                              | C||<---             D            --->|
                                                                   G  P     _  E  i  k     t  h  y  r
05 03 05 43    05 03 05 43    C3 C3 0C 43    00 63 64 82    45 0A 47 50    5F 45 69 6B    74 68 79 72


|<-- E  -->|  |<--  F  -->|  |<-- G  -->|  | H||<---            I          --->||<---   J  --->||<-- K 
                                                I  r  o     n  N  a  i     l  s 
00 00 00 00    67 00 00 00    1F 00 00 00    09 49 72 6F    6E 4E 61 69    6C 73 0A 00    00 00 00 00


 -->||<--   X   -->||<--    Y   -->|                                                              | next item -->
C8 42 02 00    00 00 03 00    00 00 00 01    00 00 00 00    00 00 00 6D    1E F7 D1 00    00 00 00 xx
```

File uses Litle Endian format, keep that in mind when reading more than one byte.

- **A.** Player Name length: Example "bjørn", due to special character has an extra byte, 0x06 --> 06 dec
- **B.** Player Name - variable in size (utf8)
- **C.** Length of the power name, example: for "GP_Eikthyr" --> 0x0A --> 10 dec 
- **D.** Name of the Power, example: "GP_Eikthyr", variable length!
- **E.** Cooldown timer of the power, microseconds? not sure yet!
- **F.** Not sure yet....
- **G.** Number of items in inventory, example: 0x1f --> 31 dec
- **H.** Length of the upcoming item name, example: "IronNails" --> 0x09 --> 09 dec
- **I.** Item name, example: "IronNails", variable length! 
- **J.** Item quantity, at the moment not sure if it's 2 bytes or 4 bytes
- **K.** no clue...

- **X.** X coordinate of item in inventory 
- **Y.** Y coordinate of item in inventory 

## Notes:
- From the end of `B` to `C` it's always 35 bytes 
- From `E[0]` to `H`, it's 12 bytes
