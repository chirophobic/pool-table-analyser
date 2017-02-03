# Pool Table Analyser

Analyse the current state of a pool table from a photo
and determine the best move/shot to make.

## Requirements

- [Go](https://golang.org/)

## Running the application
- //TODO


## Board State File

A board state file defines the current state of a board.
They are in the format:

    <board width> <board height>
    <x> <y> <type>
    <x> <y> <type>
    ...
    
Where `board width` and `board height` define the
dimensions of the board in arbitrary units.

Each following like (`<x> <y> <type>`) defines
the positions of each of the balls on the table.
`x` and `y` refer to the position of the ball on the
table. `type` refers to the type of ball (see below for
ball types).

### Ball Types

There are three ball types:

`0` This is the que ball. There should only be on of these
on the table.

`1` This is one teams ball (big/small or stripes/full).

`2` This is the other teams ball.