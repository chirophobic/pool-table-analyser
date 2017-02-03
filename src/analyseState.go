package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	"strconv"
	"image"
	"image/color"
	"image/png"
	"math"
)

const (
	BALL_TYPE_CUE = 0
	BALL_TYPE_ONE = 1
	BALL_TYPE_TWO = 2
)

type Ball struct {
	x int
	y int
	t int
}

type Hole struct {
	x int
	y int
}

type TableState struct {
	width int
	height int
	balls *[]Ball
	holes [6]Hole
}

func main() {
	args := os.Args[1:]

	if len(args) == 1 {
		run(args[0])
	} else {
		message := "Incorrect number of arguments. Expected 1 got %v\n"
		fmt.Printf(message, len(args))
	}
}

func run(stateFile string) {
	state := loadStateFromFile(stateFile)

	fmt.Printf("Width: %v, Height: %v\n", state.width, state.height)

	for _,ball := range *state.balls {
		fmt.Printf("X: %v, Y: %v, Type: %v\n", ball.x, ball.y, ball.t)
	}

	for _,hole := range state.holes {
		fmt.Printf("X: %v, Y: %v\n", hole.x, hole.y)
	}

	stateToImage(state, "output-image")
}

func loadStateFromFile(filename string) TableState {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(bytes), "\n")

	fmt.Println(string(bytes))

	tableState := TableState{}
	tableState.width, tableState.height = loadHeader(lines[0])
	tableState.balls = loadBallRows(lines[1:])
	loadHolePositions(&tableState)

	return tableState
}

func loadHeader(line string) (int, int) {
	line = strings.TrimSpace(line)
	validHeader := regexp.MustCompile(`^\d+ \d+$`)

	if !validHeader.MatchString(line) {
		panic("Invalid state header")
	}

	split := strings.Split(line, " ")
	width, err1 := strconv.Atoi(split[0])
	height, err2 := strconv.Atoi(split[1])

	panicOn(err1, err2)

	return width, height
}

func loadBallRows(lines []string) *[]Ball {
	validRow := regexp.MustCompile(`^\d+ \d+ [012]`)
	balls := make([]Ball, len(lines))

	for index,line := range lines {
		line = strings.TrimSpace(line)
		fmt.Printf("%s\n", line)

		if !validRow.MatchString(line) {
			panic("Invalid ball row line: " + line)
		}

		split := strings.Split(line, " ")
		x, err1 := strconv.Atoi(split[0])
		y, err2 := strconv.Atoi(split[1])
		t, err3 := strconv.Atoi(split[2])

		panicOn(err1, err2, err3)

		balls[index] = Ball{x, y, t}
	}

	return &balls
}

func loadHolePositions(state *TableState) {
	state.holes = [...]Hole{
		{0, 0},
		{state.width, 0},
		{0, state.height / 2},
		{state.width, state.height / 2},
		{0, state.height},
		{state.width, state.height}}
}

func stateToImage(state TableState, filepath string) {
	scale := 50
	offset := 25
	img := image.NewRGBA(image.Rect(0, 0, state.width * scale + offset * 2, state.height * scale + offset * 2))

	black := color.RGBA{0, 0, 0, 255}
	grey := color.RGBA{150, 150, 150, 255}
	white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	// Draw the table
	FillRect(img, offset, offset, state.width * scale + offset, state.height * scale + offset, green)

	// Draw the holes

	for _,hole := range state.holes {
		Circle(img, hole.x * scale + offset, hole.y * scale + offset, 10, black)
	}

	for _,ball := range *state.balls {
		var col color.Color
		if ball.t == BALL_TYPE_CUE {
			col = white
		} else if ball.t == BALL_TYPE_ONE {
			col = red
		} else {
			col = blue
		}
		Circle(img, ball.x * scale + offset, ball.y * scale + offset, 10, col)
	}

	// Get the positions of the cue and the closest ball

	var cue Ball
	for _,ball := range *state.balls {
		if ball.t == BALL_TYPE_CUE {
			cue = ball
		}
	}

	var closest Ball
	smallestDist := float64(99999999)
	for _,ball := range *state.balls {
		if ball.t == BALL_TYPE_CUE {
			continue
		}

		d := distanceBetween(cue.x, cue.y, ball.x, ball.y)
		if d < smallestDist {
			closest = ball
			smallestDist = d
		}
	}

	DiagLine(img, closest.x * scale + offset, closest.y * scale + offset, cue.x * scale + offset, cue.y * scale + offset, grey)

	f, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

// HLine draws a horizontal line
func HLine(img *image.RGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func VLine(img *image.RGBA, x, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

func DiagLine(img * image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	fx1, fy1, fx2, fy2 := float64(x1), float64(y1), float64(x2), float64(y2)
	m := (fy2 - fy1) / (fx2 - fx1)
	c := fy1 - fx1 * m

	start := math.Min(fx1, fx2)
	end := math.Max(fx1, fx2)
	fmt.Printf("%v %v\n", start, end)

	for x := int(start); x <= int(end); x++ {
		y := int(m * float64(x) + c)
		img.Set(x, y, col)
		fmt.Printf("%v %v\n", x, y)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func FillRect(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			img.Set(x, y, col)
		}
	}
}

func Rect(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	HLine(img, x1, y1, x2, col)
	HLine(img, x1, y2, x2, col)
	VLine(img, x1, y1, y2, col)
	VLine(img, x2, y1, y2, col)
}
func Circle(img *image.RGBA, x, y, radius int, col color.Color) {
	fr := float64(radius)
	for xx := -radius; xx < radius; xx++ {
		for yy := -radius; yy < radius; yy++ {
			fxx, fyy := float64(xx), float64(yy)
			d := math.Sqrt(fxx * fxx + fyy * fyy) / fr

			if d <= 1 {
				img.Set(x + xx, y + yy, col)
			}
		}
	}
}

func distanceBetween(x1, y1, x2, y2 int) float64 {
	return math.Pow(float64(x1 - x2), 2) + math.Pow(float64(y1 - y2), 2)
}

func panicOn(errors ...error) {
	for _,err := range errors {
		if err != nil {
			panic(err)
		}
	}
}
