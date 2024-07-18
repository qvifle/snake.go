package snake

import (
	"bytes"
	"fmt"
	"github.com/eiannone/keyboard"
	"log"
	"math/rand"
	"os"
	"snake/utils"
	"time"
)

type SnakePiece struct {
	x, y int
}

type Snake struct {
	area            int
	speed           time.Duration
	matrix          [][]rune
	pieces          []SnakePiece
	movement        string
	tailCoordinates SnakePiece
	bonus           SnakePiece
	status          string
}

func NewSnake() Snake {
	area, speed := getUserInput()

	bonus := getBonusCoordinates(area)
	matrix := utils.CreateMap(area)

	snake := Snake{area,
		speed,
		matrix,
		[]SnakePiece{{}},
		"right",
		SnakePiece{},
		bonus,
		"pending",
	}
	snake.render()
	return snake
}

func (s *Snake) Play() {
	ch := make(chan int)
	go func() {
		for {
			if s.status != "pending" {
				break
			}
			time.Sleep(s.speed)
			s.step()
		}
		fmt.Println(s.status)
		<-ch
		os.Exit(0)
	}()

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}

	for {
		if s.status != "pending" {
			break
		}
		char, _, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}
		s.changeMovementByChar(char)
	}

	if err := keyboard.Close(); err != nil {
		log.Fatal(err)
	}
}

// game logic
func (s *Snake) step() {
	if isWin := s.isWin(); isWin {
		s.status = "win"
		return
	}

	s.move()

	if isLose := s.isLose(); isLose {
		s.status = "lose"
		return
	}

	if snakeHead := s.pieces[len(s.pieces)-1]; s.bonus.x == snakeHead.x && s.bonus.y == snakeHead.y {
		s.levelUp()
		s.setBonus()
	}

	utils.ClearLines(s.area + 1)
	s.render()
}

func (s *Snake) move() {
	newHead := s.pieces[len(s.pieces)-1]
	switch s.movement {
	case "top":
		if newHead.y-1 < 0 {
			s.status = "lose"
			return
		}
		newHead.y = newHead.y - 1
		break
	case "right":
		if newHead.x+1 == s.area {
			s.status = "lose"
			return
		}
		newHead.x = newHead.x + 1
		break
	case "bottom":
		if newHead.y+1 == s.area {
			s.status = "lose"
			return
		}
		newHead.y = newHead.y + 1
		break
	case "left":
		if newHead.x-1 < 0 {
			s.status = "lose"
			return
		}
		newHead.x = newHead.x - 1
		break
	}

	s.tailCoordinates = s.pieces[0]
	s.pieces = append(s.pieces, newHead)
	s.pieces = s.pieces[1:]
}

func (s *Snake) isWin() bool {
	if len(s.pieces) == s.area*s.area {
		return true
	}
	return false
}

func (s *Snake) isLose() bool {
	snakeHead := s.pieces[len(s.pieces)-1]
	if len(s.pieces) == 1 {
		return false
	}

	for _, piece := range s.pieces[:len(s.pieces)-1] {
		if piece.x == snakeHead.x && piece.y == snakeHead.y {
			s.status = "lose"
			return true
		}
	}
	return false
}

func (s *Snake) setBonus() {
	if len(s.pieces) == s.area*s.area {
		return
	}
	for {
		isFound := true
		bonus := getBonusCoordinates(s.area)
		for _, piece := range s.pieces {
			if piece.x == bonus.x && piece.y == bonus.y {
				isFound = false
				break
			}
		}

		if isFound {
			s.bonus = bonus
			return
		}
	}
}

func getBonusCoordinates(area int) SnakePiece {
	return SnakePiece{rand.Intn(area), rand.Intn(area)}
}

func (s *Snake) levelUp() {
	s.pieces = append([]SnakePiece{s.tailCoordinates}, s.pieces...)
}

func (s *Snake) changeMovementByChar(char rune) {
	switch char {
	case 'w', 'ц':
		s.setMovement("top")
	case 'd', 'в':
		s.setMovement("right")
	case 's', 'ы':
		s.setMovement("bottom")
	case 'a', 'ф':
		s.setMovement("left")
	}
}

// render section
func (s *Snake) render() {

	s.matrix = utils.CreateMap(s.area)
	s.matrix[s.bonus.y][s.bonus.x] = utils.BonusSymbol

	for _, piece := range s.pieces {
		s.matrix[piece.y][piece.x] = utils.SquareSymbol
	}

	fmt.Printf("Score: %d\n", len(s.pieces))
	s.print()
}

func (s *Snake) print() {
	for i := 0; i < s.area; i++ {
		var row bytes.Buffer
		for j := 0; j < s.area; j++ {
			row.WriteString(fmt.Sprintf(" %c ", s.matrix[i][j]))
		}
		fmt.Println(row.String())
	}
}

// etc
func getUserInput() (int, time.Duration) {
	var area int
	var speed int

	for {
		fmt.Println("Input area (even int): ")
		_, err := fmt.Scan(&area)

		if err != nil {
			continue
		}

		utils.ClearLines(2)
		if area > 1 && area%2 == 0 {
			break
		}
	}

	for {
		fmt.Printf("Select game mode:\n1.Easy  2.Normal  3.Hard  4.Expert\n")
		_, err := fmt.Scan(&speed)

		if err != nil {
			continue
		}
		utils.ClearLines(3)
		if speed > 0 && speed < 5 {
			break
		}

	}

	gameSpeed := time.Second / time.Duration(speed)

	return area, gameSpeed
}

func (s *Snake) isBackMove(movement string) {
	if s.movement == movement {
		s.status = "lose"
		return
	}
}

func (s *Snake) setMovement(movement string) {
	if (s.movement == "top" && movement == "bottom" ||
		s.movement == "bottom" && movement == "top" ||
		s.movement == "right" && movement == "left" ||
		s.movement == "left" && movement == "right") &&
		len(s.pieces) > 1 {
		return
	}

	s.movement = movement
}
