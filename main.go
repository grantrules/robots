package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/mattn/go-tty"
)

func main() {

	// handle CTRL C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	width, height := getTerminalSize()

	if width == 0 || height == 0 {
		exitWithMessage("Terminal size is too small")
	}

	// hide the cursor
	hideCursor()

	game := NewGame(width, height)

	game.Start()

	go func() {
		for range c {
			game.over()
		}
	}()

	go game.handleInput()
	for {
		game.draw()
		time.Sleep(time.Millisecond * 50)
	}

}

func distance(a, b position) int {
	return abs(a[0]-b[0]) + abs(a[1]-b[1])
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func exitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(0)
}

func (g *Game) over() {
	// clear
	showCursor()
	clear()
	exitWithMessage("Game over")
}

func (g *Game) handleInput() {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		char, err := tty.ReadRune()
		if err != nil {
			panic(err)
		}

		moved := true
		switch char {
		case 'w':
			g.movePlayer(0, -1)
		case 'a':
			g.movePlayer(-1, 0)
		case 's':
			g.movePlayer(0, 1)
		case 'd':
			g.movePlayer(1, 0)
		case ' ':
		default:
			moved = false
		}
		if moved {
			g.moveRobots()
			g.checkRobotCollisions()
			if g.checkWinCondition() {
				g.dialog = makeDialogFunc(10, 10, 10, 10, "You win!")
				break
			}
		}
	}
}

func (g *Game) checkIfPlayerCanMove(x, y int) bool {
	newPlayerPos := position{g.player[0] + x, g.player[1] + y}

	if newPlayerPos[0] < 0 || newPlayerPos[0] >= g.Width {
		return false
	}
	if newPlayerPos[1] < 0 || newPlayerPos[1] >= g.Height {
		return false
	}

	for i := range g.robots {
		r := &g.robots[i]
		rX, rY := r.getRobotMove(newPlayerPos)
		newPos := position{r.position[0] + rX, r.position[1] + rY}

		if (newPos == position{x, y}) {
			return false
		}
	}

	return true
}

func (g *Game) movePlayer(x, y int) {
	if g.checkIfPlayerCanMove(x, y) {
		g.player[0] += x
		g.player[1] += y
	}
}

func (g *Game) moveRobots() {
	for i := range g.robots {
		r := &g.robots[i]
		if r.alive {
			r.move(r.getRobotMove(g.player))
		}
	}
}

func (g *Game) checkRobotCollisions() {
	for i := range g.robots {
		for j := range g.robots {
			if i != j && g.robots[i].position == g.robots[j].position {
				g.robots[i].alive = false
				g.robots[j].alive = false
			}
		}
	}
}

func (g *Game) checkWinCondition() bool {
	for i := range g.robots {
		if g.robots[i].alive {
			return false
		}
	}
	return true
}

func (g *Game) isRobotInPosition(pos position) bool {
	for i := range g.robots {
		r := g.robots[i]
		x, y := r.getRobotMove(g.player)

		newPos := position{r.position[0] + x, r.position[1] + y}
		if newPos == pos {
			return true
		}
	}
	return false
}

// Start begins the game
func (g *Game) Start() {

	g.player = position{rand.Intn(g.Width), rand.Intn(g.Height)}

	for i := 0; i < 25; i++ {
		robot := g.NewRobot()
		for distance(g.player, robot.position) < 5 {
			robot = g.NewRobot()
		}
		g.robots = append(g.robots, robot)
	}

}

func (g *Game) NewRobot() Robot {
	robot := Robot{
		position: position{rand.Intn(g.Width), rand.Intn(g.Height)},
		alive:    true,
	}
	return robot
}

type Game struct {
	Width, Height int

	player position

	robots []Robot

	dialog func()
}

type position [2]int

func NewGame(width, height int) Game {
	game := Game{
		Width:  width,
		Height: height,
	}
	return game
}

func (g *Game) draw() {
	clear()

	for i := range g.robots {
		r := g.robots[i]
		icon := "🤖"
		if !r.alive {
			icon = "🗑️"
		}
		x, y := r.position[0], r.position[1]
		printAt(x, y, icon)
	}

	printAt(g.player[0], g.player[1], "👨")

	if g.dialog != nil {
		g.dialog()
	}

}

func makeDialogFunc(x, y, width, height int, title string) func() {
	return func() {
		dialog(x, y, width, height, title)
	}
}
