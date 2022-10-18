package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	cursor "github.com/ahmetalpbalkan/go-cursor"
	keyboard "github.com/eiannone/keyboard"
)

var _WIDTH = 64
var _HEIGHT = 32

type Position struct {
	x int
	y int
}

type Direction struct {
	x int
	y int
}

func Repeat(x int, y int) int {
	if x < 0 {
		return y
	} else if x > y {
		return 0
	}
	return x
}

func readKeyboard(dir *Direction, lastDir *Direction, snakeHead *string) {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Print(cursor.MoveTo(_HEIGHT+5, 0), "Press ESC or Q to quit")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		if event.Key == keyboard.KeyArrowDown && lastDir.y == 0 {
			dir.x = 0
			dir.y = 1
			*snakeHead = "👇"
		}
		if event.Key == keyboard.KeyArrowUp && lastDir.y == 0 {
			dir.x = 0
			dir.y = -1
			*snakeHead = "👆"
		}
		if event.Key == keyboard.KeyArrowLeft && lastDir.x == 0 {
			dir.y = 0
			dir.x = -1
			*snakeHead = "👈"
		}
		if event.Key == keyboard.KeyArrowRight && lastDir.x == 0 {
			dir.y = 0
			dir.x = 1
			*snakeHead = "👉"
		}
		if event.Key == keyboard.KeyEsc || event.Rune == 'q' || event.Rune == 'Q' {
			fmt.Println(cursor.MoveTo(_HEIGHT+3, 0))
			fmt.Print(cursor.Show())
			os.Exit(0)
		}
	}
}

func drawBox(width int, height int) {
	fmt.Print(cursor.MoveTo(0, 0))
	fmt.Println("┌" + strings.Repeat("─", (width*2)+2) + "┐")
	for i := 0; i < height; i++ {
		fmt.Println("│ " + strings.Repeat(" ", (width*2)) + " │")
	}
	fmt.Println("└" + strings.Repeat("─", (width*2)+2) + "┘")
}

var emojis [90]string = [90]string{
	"👆", "😂", "😝", "😁", "😱", "👉", "🙌", "🍻", "🔥", "🌈", "🎈", "🌹", "💄", "🎀", "⚽", "🎾", "🏁", "😡", "👿", "🐻", "🐶", "🐬", "🐟", "🍀", "👀", "🚗", "🍎", "💝", "💙", "👌", "😍", "😉", "😓", "😳", "💪", "💩", "🍸", "🔑", "💖", "🌟", "🎉", "🌺", "🎶", "👠", "🏈", "⚾", "🏆", "👽", "💀", "🐵", "🐮", "🐩", "🐎", "💣", "👃", "👂", "🍓", "💘", "💜", "👊", "💋", "😘", "😜", "😵", "🙏", "👋", "🚽", "💃", "💎", "🚀", "🌙", "🎁", "⛄", "🌊", "⛵", "🏀", "🎱", "💰", "👶", "👸", "🐰", "🐷", "🐍", "🐫", "🔫", "👄", "🚲", "🍉", "💛", "💚",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Print(cursor.Hide())
	fmt.Print(cursor.ClearEntireScreen())

	_WIDTH, _HEIGHT, _ = term.GetSize(int(os.Stdin.Fd()))

	_WIDTH = _WIDTH/2 - 2
	_HEIGHT -= 6
	//_WIDTH = 32
	//_HEIGHT = 32

	playerPosition := Position{0, 0}
	foodPosition := Position{rand.Int() % _WIDTH, rand.Int() % _HEIGHT}
	direction := Direction{1, 0}
	lastDirection := Direction{1, 0}
	score := 24
	snakeHead := "👉"

	field := make([][]int, _HEIGHT)
	for i := 0; i < _HEIGHT; i++ {
		field[i] = make([]int, _WIDTH)
	}
	needsUpdate := make([][]bool, _HEIGHT)
	for i := 0; i < _HEIGHT; i++ {
		needsUpdate[i] = make([]bool, _WIDTH)
		for j := 0; j < _WIDTH; j++ {
			needsUpdate[i][j] = true
		}
	}

	field[playerPosition.y][playerPosition.x] = score
	field[foodPosition.y][foodPosition.x] = -1

	go readKeyboard(&direction, &lastDirection, &snakeHead)

	drawBox(_WIDTH, _HEIGHT)

	var lastRun int64 = 0

	for {
		currentTime := time.Now().UnixMilli()
		if currentTime-lastRun < 68 {
			continue
		}
		lastRun = currentTime
		if playerPosition == foodPosition {
			score++
			if score == _WIDTH*_HEIGHT {
				fmt.Println(cursor.MoveTo(_HEIGHT/2+2, _WIDTH) + cursor.Esc + "[42m" + cursor.Esc + "[30m  YOU WIN  " + cursor.Esc + "[0m")
				fmt.Print(cursor.MoveTo(_HEIGHT+4, 0))
				os.Exit(0)
			}
			field[foodPosition.y][foodPosition.x] = score
			foodPosition = Position{rand.Int() % _WIDTH, rand.Int() % _HEIGHT}
			field[foodPosition.y][foodPosition.x] = -1
			needsUpdate[foodPosition.y][foodPosition.x] = true
		}
		fmt.Println(cursor.MoveTo(_HEIGHT+3, 0)+"Score:", score)

		for h := 0; h < _HEIGHT; h++ {
			for w := 0; w < _WIDTH; w++ {
				if field[h][w] > 0 {
					fmt.Print(cursor.MoveTo(h+2, w+w+3), emojis[field[h][w]%90])
					field[h][w]--
					needsUpdate[h][w] = true
				} else if field[h][w] == -1 && needsUpdate[h][w] {
					fmt.Print(cursor.MoveTo(h+2, w+w+3), "🍖")
					needsUpdate[h][w] = false
				} else if needsUpdate[h][w] {
					fmt.Print(cursor.MoveTo(h+2, w+w+3), "•")
					needsUpdate[h][w] = false
				}
			}
		}

		lastDirection = direction
		playerPosition.x += direction.x
		playerPosition.y += direction.y

		if playerPosition.x < 0 || playerPosition.x >= _WIDTH || playerPosition.y < 0 || playerPosition.y >= _HEIGHT || field[playerPosition.y][playerPosition.x] > 0 {
			fmt.Print(cursor.MoveTo(playerPosition.y+2, playerPosition.x+playerPosition.x+3) + cursor.Esc + "[31m❌" + cursor.Esc + "[0m")
			fmt.Print(cursor.MoveTo(_HEIGHT/2+2, _WIDTH-3) + cursor.Esc + "[41m  YOU LOSE  " + cursor.Esc + "[0m")
			fmt.Print(cursor.MoveTo(_HEIGHT+4, 0))
			os.Exit(0)
		}

		field[playerPosition.y][playerPosition.x] = score
	}
}
