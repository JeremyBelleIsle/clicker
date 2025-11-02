package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Structure d’un sac d’argent dans le fond
type MoneyBagBG struct {
	x, y  float64
	speed float64
	scale float64
	angle float64
}

type Game struct {
	Clicks              int
	ClickPower          int
	AutoClickPower      int
	framesCNT           int
	ChanceFramesCNTAnim int
	clicPrecedent       bool
	MoneyRequireC       int
	MoneyRequireA       int
	clickAnim           float64
	MoneyBag            *ebiten.Image
	BagsBG              []MoneyBagBG
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	rand.Seed(time.Now().UnixNano())
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func Within(x, y, rx, ry, rw, rh int) bool {
	return x >= rx && x <= rx+rw && y >= ry && y <= ry+rh
}

func (g *Game) spawnBackgroundBags() {
	// Crée un fond rempli de sacs à des positions aléatoires
	for i := 0; i < 25; i++ {
		g.BagsBG = append(g.BagsBG, MoneyBagBG{
			x:     rand.Float64() * 640,
			y:     -200 + rand.Float64()*480,
			speed: 0.5 + rand.Float64()*1.5,
			scale: 0.1 + rand.Float64()*0.2,
			angle: rand.Float64() * 360,
		})
	}
}

func (g *Game) Update() error {
	if g.ChanceFramesCNTAnim > 0 {
		g.ChanceFramesCNTAnim--
	}
	g.MoneyRequireC = g.ClickPower * 30
	g.MoneyRequireA = (g.AutoClickPower + 1) * 50

	// Auto clicks chaque seconde
	g.framesCNT++
	if g.framesCNT >= 60 {
		g.Clicks += g.AutoClickPower
		g.framesCNT = 0
	}

	// 💸 défilement des sacs d’arrière-plan
	for i := range g.BagsBG {
		g.BagsBG[i].x += g.BagsBG[i].speed
		g.BagsBG[i].y += g.BagsBG[i].speed * 0.8

		if g.BagsBG[i].x > 700 || g.BagsBG[i].y > 520 {
			g.BagsBG[i].x = rand.Float64()*-100 - 50
			g.BagsBG[i].y = -200 + rand.Float64()*480 - 50
			g.BagsBG[i].speed = 0.5 + rand.Float64()*1.5
		}
	}

	clicActuel := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	x, y := ebiten.CursorPosition()

	// 💥 Clic sur le sac principal
	if clicActuel && !g.clicPrecedent {
		imgW, imgH := g.MoneyBag.Size()
		imgX := 320 - int(float64(imgW)*0.6/2)
		imgY := 240 - int(float64(imgH)*0.6/2)
		if Within(x, y, imgX, imgY, int(float64(imgW)*0.6), int(float64(imgH)*0.6)) {
			chance := rand.Intn(1000) + 1 // entre 1 et 1000
			switch {
			case chance <= 5:
				g.Clicks += g.ClickPower * 100
				g.ChanceFramesCNTAnim = 120
				fmt.Println("💎 MYTHIQUE JACKPOT !!!")
			case chance <= 25:
				g.Clicks += g.ClickPower * 20
				fmt.Println("💰 Légendaire !")
			case chance <= 100:
				g.Clicks += g.ClickPower * 5
				fmt.Println("⭐ Coup de chance !")
			default:
				g.Clicks += g.ClickPower
			}

			g.clickAnim = 0
		}
	}

	// ⚙️ Upgrade click
	if clicActuel && Within(x, y, 450, 50, 175, 50) && g.Clicks >= g.MoneyRequireC {
		g.Clicks -= g.MoneyRequireC
		g.ClickPower++
	}

	// ⚙️ Upgrade auto-click
	if clicActuel && Within(x, y, 450, 150, 175, 50) && g.Clicks >= g.MoneyRequireA {
		g.Clicks -= g.MoneyRequireA
		g.AutoClickPower++
	}

	g.clicPrecedent = clicActuel

	// Animation du clic
	if g.clickAnim < 0.3 {
		g.clickAnim += 1.0 / 60.0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// couleur de derrière
	ebitenutil.DrawRect(screen, 0, 0, 640, 480, color.RGBA{255, 182, 193, 40})
	ebitenutil.DrawRect(screen, 0, 0, 640, 480, color.RGBA{0, 130, 0, 40})
	// 🌌 Dessin du fond animé
	for _, b := range g.BagsBG {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-32, -32)
		op.GeoM.Scale(b.scale, b.scale)
		op.GeoM.Translate(b.x, b.y)
		op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 150}) // transparence légère
		screen.DrawImage(g.MoneyBag, op)
	}

	imgW, imgH := g.MoneyBag.Size()
	scale := 0.60
	if g.clickAnim < 0.3 {
		scale = 0.60 - 0.06*math.Sin((g.clickAnim/0.3)*math.Pi)
	}

	// 💰 Sac principal au centre
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(imgW)/2, -float64(imgH)/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(320, 240)
	screen.DrawImage(g.MoneyBag, op)

	// 💵 Affichage argent
	opMoney := &text.DrawOptions{}
	opMoney.GeoM.Translate(20, 30)
	opMoney.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 255})
	text.Draw(screen, fmt.Sprintf("Money: %d", g.Clicks), &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   14,
	}, opMoney)

	// 🔹 Upgrade click
	ebitenutil.DrawRect(screen, 450, 50, 179, 59, color.RGBA{22, 13, 200, 255})
	ebitenutil.DrawRect(screen, 440, 40, 200, 79, color.RGBA{22, 13, 200, 100})
	opClick := &text.DrawOptions{}
	opClick.GeoM.Translate(460, 70)
	opClick.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
	text.Draw(screen, "Upgrade click", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   12,
	}, opClick)

	opClickInfo := &text.DrawOptions{}
	opClickInfo.GeoM.Translate(460, 100)
	opClickInfo.ColorScale.ScaleWithColor(color.RGBA{255, 255, 0, 255})
	text.Draw(screen,
		fmt.Sprintf("Power: %d | Cost: %d", g.ClickPower, g.MoneyRequireC),
		&text.GoTextFace{
			Source: mplusFaceSource,
			Size:   10,
		}, opClickInfo)

	// 🔹 Upgrade auto click
	ebitenutil.DrawRect(screen, 450, 150, 179, 59, color.RGBA{22, 13, 200, 255})
	ebitenutil.DrawRect(screen, 440, 140, 200, 79, color.RGBA{22, 13, 200, 100})
	opAuto := &text.DrawOptions{}
	opAuto.GeoM.Translate(460, 170)
	opAuto.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
	text.Draw(screen, "Upgrade auto click", &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   10,
	}, opAuto)

	opAutoInfo := &text.DrawOptions{}
	opAutoInfo.GeoM.Translate(460, 190)
	opAutoInfo.ColorScale.ScaleWithColor(color.RGBA{255, 255, 0, 255})
	text.Draw(screen,
		fmt.Sprintf("Power: %d | Cost: %d", g.AutoClickPower, g.MoneyRequireA),
		&text.GoTextFace{
			Source: mplusFaceSource,
			Size:   10,
		}, opAutoInfo)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("💸 Money Clicker Deluxe")

	game := &Game{
		ClickPower:    1,
		MoneyRequireC: 30,
		MoneyRequireA: 50,
	}

	var err error
	game.MoneyBag, _, err = ebitenutil.NewImageFromFile("money bag.png")
	if err != nil {
		log.Fatal(err)
	}

	game.spawnBackgroundBags()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
