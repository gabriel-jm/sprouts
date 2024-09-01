package internal

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	Sprite rl.Texture2D
	Source rl.Rectangle
	Dest   rl.Rectangle
	Speed  float32
	Dir    int
	Frame  int
}

func NewPlayer() Player {
	return Player{
		Sprite: rl.LoadTexture("res/Characters/basic_char.png"),
		Source: rl.NewRectangle(0, 0, 48, 48),
		Dest:   rl.NewRectangle(200, 200, 60, 60),
		Speed:  1.4,
	}
}

func (p *Player) ResetSource() {
	p.Source.X = p.Source.Width * float32(p.Frame)
}

func (p *Player) Move(isMoving bool) {
	if !isMoving && p.Frame > 1 {
		p.Frame = 0
	}

	if p.Frame > 3 {
		p.Frame = 0
	}

	p.Source.X = p.Source.Width * float32(p.Frame)
	p.Source.Y = p.Source.Height * float32(max(0, p.Dir))
}
