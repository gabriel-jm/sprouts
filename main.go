package main

import (
	"fmt"
	"os"
	"sprouts/internal"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1000
	screenHeight = 480
)

var (
	gameRunning     = true
	backgroundColor = rl.NewColor(147, 211, 196, 255)

	grassSprite  rl.Texture2D
	hillSprite   rl.Texture2D
	fenceSprite  rl.Texture2D
	houseSprite  rl.Texture2D
	waterSprite  rl.Texture2D
	tilledSprite rl.Texture2D

	player internal.Player

	tileDest  rl.Rectangle
	tileSrc   rl.Rectangle
	tileMap   []int
	srcMap    []string
	mapWidth  int
	mapHeight int

	frameCount int

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

func update() {
	player.ResetSource()

	isMoving := input()
	gameRunning = !rl.WindowShouldClose()

	frameCount++
	player.Move(isMoving)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(
		float32(player.Dest.X-(player.Dest.Width/2)),
		float32(player.Dest.Y-(player.Dest.Height/2)),
	)
}

func input() bool {
	var speedX float32 = 0.0
	var speedY float32 = 0.0

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		speedY = -player.Speed
		player.Dir = 1
	}

	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		speedY = player.Speed
		player.Dir = 0
	}

	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		speedX = -player.Speed
		player.Dir = 2
	}

	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		speedX = player.Speed
		player.Dir = 3
	}

	if speedX != 0.0 || speedY != 0.0 {
		player.Dest.X += float32(speedX)
		player.Dest.Y += float32(speedY)
		if frameCount%8 == 1 {
			player.Frame++
		}

		return true
	}

	if frameCount%45 == 1 {
		player.Frame++
	}

	return false
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(backgroundColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func drawScene() {
	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i%mapWidth)
			tileDest.Y = tileDest.Height * float32(i/mapWidth)

			var tex rl.Texture2D

			if srcMap[i] == "g" {
				tex = grassSprite
			}
			if srcMap[i] == "l" {
				tex = hillSprite
			}
			if srcMap[i] == "f" {
				tex = fenceSprite
			}
			if srcMap[i] == "h" {
				tex = houseSprite
			}
			if srcMap[i] == "w" {
				tex = waterSprite
			}
			if srcMap[i] == "t" {
				tex = tilledSprite
			}

			if srcMap[i] == "h" || srcMap[i] == "f" {
				tileSrc.X = 64
				tileSrc.Y = 64
				rl.DrawTexturePro(
					grassSprite,
					tileSrc,
					tileDest,
					rl.NewVector2(tileDest.Width, tileDest.Height),
					0,
					rl.White,
				)
			}

			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))
			rl.DrawTexturePro(
				tex,
				tileSrc,
				tileDest,
				rl.NewVector2(tileDest.Width, tileDest.Height),
				0,
				rl.White,
			)
		}
	}

	rl.DrawTexturePro(
		player.Sprite,
		player.Source,
		player.Dest,
		rl.NewVector2(player.Dest.Width, player.Dest.Height),
		0,
		rl.White,
	)
}

func loadMap(fileName string) {
	file, err := os.ReadFile("maps/" + fileName)

	if err != nil {
		fmt.Println(file)
		os.Exit(1)
	}

	tileMapList := strings.Split(strings.Replace(string(file), "\n", " ", -1), " ")
	mapWidth = -1
	mapHeight = -1

	for i := 0; i < len(tileMapList); i++ {
		tileNumber, _ := strconv.ParseInt(tileMapList[i], 10, 64)

		if mapWidth == -1 {
			mapWidth = int(tileNumber)
		} else if mapHeight == -1 {
			mapHeight = int(tileNumber)
		} else if i < mapWidth*mapHeight+2 {
			tileMap = append(tileMap, int(tileNumber))
		} else {
			srcMap = append(srcMap, tileMapList[i])
		}
	}
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "Sprouts")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("res/Tilesets/Grass.png")
	hillSprite = rl.LoadTexture("res/Tilesets/Hills.png")
	fenceSprite = rl.LoadTexture("res/Tilesets/Fences.png")
	houseSprite = rl.LoadTexture("res/Tilesets/Wooden_House_Walls_Tilset.png")
	waterSprite = rl.LoadTexture("res/Tilesets/Water.png")
	tilledSprite = rl.LoadTexture("res/Tilesets/Tilled_Dirt.png")

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	player = internal.NewPlayer()

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/music/hopeful.mp3")
	musicPaused = false
	rl.SetMusicVolume(music, 0.5)
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(player.Dest.X-(player.Dest.Width/2)), float32(player.Dest.Y-(player.Dest.Height/2))),
		0.0,
		1.0,
	)

	cam.Zoom = 3

	loadMap("one.map")
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(player.Sprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {
	defer quit()

	for gameRunning {
		update()
		render()
	}
}
