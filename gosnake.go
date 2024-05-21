package main

import (
	"log"
	"time"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type command int;

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
	EXIT
)

func game_loop(rune_ch <-chan command, s *tcell.Screen) int {
	score := 0;
	snake_x := make([]int, 2);
	snake_y := make([]int, 2);
	snake_dir := DOWN;
	food_x := 5;
	food_y := 5;

	snake_y[0] = 1;
	snake_y[1] = 0;

	sa := *s;
	w, h := sa.Size();

	for {
		sa.Show();
		sa.Clear();

		x, y := snake_x[0], snake_y[0];

		//Check for new input
		select {
		case cmd := <-rune_ch:
			switch cmd {
			case EXIT:
				sa.Clear();
				sa.Show();
				return 0;
			case UP:
				if snake_dir > 0x01 {
					snake_dir = UP;
				}
			case DOWN:
				if snake_dir > 0x01 {
					snake_dir = DOWN;
				}
			case LEFT:
				if snake_dir < 0x02 {
					snake_dir = LEFT;
				}
			case RIGHT:
				if snake_dir < 0x02 {
					snake_dir = RIGHT;
				}
			default:
			}
		case <-time.After(64 * time.Millisecond):
		}

		if snake_dir == UP {
			y--;
		} else if snake_dir == DOWN {
			y++;
		} else if snake_dir == LEFT {
			x--;
		} else if snake_dir == RIGHT {
			x++;
		}

		//Check if out of bounds
		if x < 0 || x > w-1 || y < 0 || y > h-1 {
			time.Sleep(1 * time.Second);
			return score;
		}

		//Check if eating self
		for i := 0; i < len(snake_x); i++ {
			if snake_x[i] == x {
				if snake_y[i] == y {
					time.Sleep(1 * time.Second);
					return score;
				}
			}
		}
		
		//Check if got food
		if x != food_x || y != food_y {
			//Pop last element
			_, snake_x = snake_x[len(snake_x)-1], snake_x[:len(snake_x)-1];
			_, snake_y = snake_y[len(snake_y)-1], snake_y[:len(snake_y)-1];
		} else {
			//Move food somewhere else
			score++;
			food_x, food_y = rand.Intn(w), rand.Intn(h);
		}

		//Push Front
		snake_x = append([]int{x}, snake_x...);	
		snake_y = append([]int{y}, snake_y...);	
		
		//Draw food
		sa.SetContent(food_x, food_y, 'f', nil, tcell.StyleDefault);

		//Draw snake
		for i := 0; i < len(snake_x); i++ {
			sa.SetContent(snake_x[i], snake_y[i], '8', nil, tcell.StyleDefault);
		}
	}
}

func rune_listener(rune_ch chan <- command, s *tcell.Screen) {
	for {
		sa := *s;
		k := sa.PollEvent();
		switch k := k.(type) {
		case *tcell.EventResize:
			sa.Sync();
		case *tcell.EventKey:
			switch k.Key() {
			case tcell.KeyEscape:
				rune_ch<- EXIT;
				return;
			case tcell.KeyCtrlC:
				rune_ch<- EXIT;
				return;
			case tcell.KeyUp:
				rune_ch<- UP;
			case tcell.KeyDown:
				rune_ch<- DOWN;
			case tcell.KeyLeft:
				rune_ch<- LEFT;
			case tcell.KeyRight:
				rune_ch<- RIGHT;
			}
		}
	}
}

func main() {
	s, err := tcell.NewScreen();
	if err != nil {
		log.Fatal(err);
	}
	if err := s.Init(); err != nil {
		log.Fatal(err);
	}

	quit := func() {
		maybePanic := recover();
		s.Fini();
		if maybePanic != nil {
			panic(maybePanic);
		}			
	}
	defer quit();

	rune_ch := make(chan command);
	go rune_listener(rune_ch, &s);
	game_loop(rune_ch, &s);	
	return;
}
