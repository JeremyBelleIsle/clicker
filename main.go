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
	BestClicks          int
	ClickPower          int
	AutoClickPower      int
	framesCNT           int
	ChanceFramesCNTAnim int
	chanceAnim          int
	gauge               int
	PrestigePoints      int
	Âmes                int
	Store               bool
	clicPrecedent       bool
	MoneyRequireC       int
	MoneyRequireA       int
	phrase              string
	clickAnim           float64
	MoneyBag            *ebiten.Image
	BagsBG              []MoneyBagBG
	RebirthTF           bool
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
	if !g.Store {
		if g.Clicks > g.BestClicks {
			g.BestClicks = g.Clicks
		}
		if g.gauge > 0 {
			g.gauge--
		} else {
			g.gauge = 0
		}
		if g.ChanceFramesCNTAnim > 0 {
			g.ChanceFramesCNTAnim--
		}
		g.MoneyRequireC = g.ClickPower * 30
		g.MoneyRequireA = (g.AutoClickPower + 1) * 50

		// Auto clicks chaque seconde
		g.framesCNT++
		if g.framesCNT >= 60 {
			g.Clicks += g.AutoClickPower * (g.PrestigePoints + 3)
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
				chance := rand.Intn(5000) + 1 // entre 1 et 1000
				switch {
				case chance <= 25:
					if g.gauge >= 200 {
						g.Clicks += ((g.ClickPower * 100) * 2) * (g.PrestigePoints + 1)
					} else {
						g.Clicks += (g.ClickPower * 100) * (g.PrestigePoints + 1)
					}
					g.ChanceFramesCNTAnim = 300
					fmt.Println("💎 MYTHIQUE JACKPOT !!!")
					g.chanceAnim = 3
				case chance <= 50:
					if g.gauge >= 200 {
						g.Clicks += ((g.ClickPower * 20) * 2) * (g.PrestigePoints + 1)
					} else {
						g.Clicks += (g.ClickPower * 20) * (g.PrestigePoints + 1)
					}
					g.ChanceFramesCNTAnim = 120
					fmt.Println("💰 Légendaire !")
					g.chanceAnim = 2
				case chance <= 100:
					if g.gauge >= 200 {
						g.Clicks += ((g.ClickPower * 5) * 2) * (g.PrestigePoints + 1)
					} else {
						g.Clicks += (g.ClickPower * 5) * (g.PrestigePoints + 1)
					}
					g.ChanceFramesCNTAnim = 50
					fmt.Println("⭐ Coup de chance !")
					g.chanceAnim = 1
				default:
					if g.gauge >= 200 {
						g.Clicks += (g.ClickPower * 2) * (g.PrestigePoints + 1)
					} else {
						g.Clicks += g.ClickPower * (g.PrestigePoints + 1)
					}
					if g.gauge < 260 {
						g.gauge += 15
					} else {
						g.gauge = 260
					}
				}
				g.clickAnim = 0
			}
			if Within(x, y, 10, 100, 150, 35) && g.RebirthTF {
				g.PrestigePoints++
				g.Âmes += rand.Intn(5) + 2
				g.Clicks = 0
				g.ClickPower = 1
				g.AutoClickPower = 0
				g.RebirthTF = false
				g.gauge = 0
				g.ChanceFramesCNTAnim = 180
				g.chanceAnim = 3
				fmt.Println("🌟 REBIRTH ACTIVÉ ! +1 Prestige")
			}
		}
		if g.Clicks >= 1_000_000 {
			g.RebirthTF = true
		}

		// ⚙️ Upgrade click
		if clicActuel && Within(x, y, 450, 50, 175, 50) && g.Clicks >= g.MoneyRequireC && !g.clicPrecedent {
			g.Clicks -= g.MoneyRequireC
			g.ClickPower++
		}

		// ⚙️ Upgrade auto-click
		if clicActuel && Within(x, y, 450, 150, 175, 50) && g.Clicks >= g.MoneyRequireA && !g.clicPrecedent {
			g.Clicks -= g.MoneyRequireA
			g.AutoClickPower++
		}

		g.clicPrecedent = clicActuel

		// Animation du clic
		if g.clickAnim < 0.3 {
			g.clickAnim += 1.0 / 60.0
		}
		if Within(x, y, 25, 70, 190, 50) && clicActuel {
			g.Store = true
		}
	} else {
		clicActuel := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		x, y := ebiten.CursorPosition()
		if Within(x, y, 260, 400, 120, 30) && clicActuel && !g.clicPrecedent {
			g.Store = false
		}
		g.clicPrecedent = clicActuel
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.Store {
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
		if g.RebirthTF {
			ebitenutil.DrawRect(screen, 10, 100, 150, 35, color.RGBA{200, 200, 200, 120})
			opRebirth := &text.DrawOptions{}
			opRebirth.GeoM.Translate(10, 110)
			opRebirth.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 255})
			text.Draw(screen, "Rebirth", &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   20,
			}, opRebirth)
		}
		maxGauge := 260.0
		gaugeWidth := (float64(g.gauge) / maxGauge) * 300.0

		ebitenutil.DrawRect(screen, 20, 430, 300, 20, color.RGBA{200, 200, 200, 120})
		var gaugeColor color.RGBA
		switch {
		case g.gauge < 80:
			gaugeColor = color.RGBA{135, 206, 250, 255}
		case g.gauge < 180:
			gaugeColor = color.RGBA{0, 191, 255, 255}
		default:
			gaugeColor = color.RGBA{255, 105, 180, 255}
		}

		ebitenutil.DrawRect(screen, 20, 430, gaugeWidth, 20, gaugeColor)
		// contour
		ebitenutil.DrawRect(screen, 20, 429, 300, 1, color.RGBA{255, 255, 255, 150})
		ebitenutil.DrawRect(screen, 20, 450, 300, 1, color.RGBA{255, 255, 255, 150})

		// 💵 Affichage argent
		opMoney := &text.DrawOptions{}
		opMoney.GeoM.Translate(20, 30)
		opMoney.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 255})
		text.Draw(screen, fmt.Sprintf("Money: %d", g.Clicks), &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   14,
		}, opMoney)
		ebitenutil.DrawRect(screen, 20, 50, 190, 50, color.RGBA{22, 13, 200, 255})
		opRebirthStore := &text.DrawOptions{}
		opRebirthStore.GeoM.Translate(25, 70)
		opRebirthStore.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 255})
		text.Draw(screen, "Rebirth Store", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   14,
		}, opRebirthStore)
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
		if g.ChanceFramesCNTAnim > 0 {
			switch g.chanceAnim {
			case 1:
				opchance := &text.DrawOptions{}
				opchance.GeoM.Translate(0, 225)
				opchance.ColorScale.ScaleWithColor(color.RGBA{135, 206, 250, 255})
				text.Draw(screen, "Coup de chance !", &text.GoTextFace{
					Source: mplusFaceSource,
					Size:   45,
				}, opchance)

			case 2:
				opchance := &text.DrawOptions{}
				opchance.GeoM.Translate(0, 225)
				opchance.ColorScale.ScaleWithColor(color.RGBA{135, 206, 250, 255})
				text.Draw(screen, "Légendaire !", &text.GoTextFace{
					Source: mplusFaceSource,
					Size:   45,
				}, opchance)
			case 3:
				opchance := &text.DrawOptions{}
				opchance.GeoM.Translate(0, 225)
				opchance.ColorScale.ScaleWithColor(color.RGBA{135, 206, 250, 255})
				text.Draw(screen, "MYTHIQUE JACKPOT !!!", &text.GoTextFace{
					Source: mplusFaceSource,
					Size:   45,
				}, opchance)
			}
		}
		switch {
		case g.BestClicks < 100:
			g.phrase = "Ton empire grandit… un jour, tu atteindras le million."
		case g.BestClicks < 1000:
			g.phrase = "Tu avances vers le seuil du vrai pouvoir."
		case g.BestClicks < 5000:
			g.phrase = "Chaque clic t’approche d’un changement irréversible."
		case g.BestClicks < 15000:
			g.phrase = "À un million, une nouvelle vie t’attend."
		case g.BestClicks < 30000:
			g.phrase = "Ta fortune approche des limites humaines…"
		case g.BestClicks < 50000:
			g.phrase = "Tu ressens l’appel du Rebirth, n’est-ce pas ?"
		case g.BestClicks < 100000:
			g.phrase = "Le million n’est plus un rêve… il est à portée de main."
		case g.BestClicks < 500000:
			g.phrase = "Chaque clic résonne plus fort… tu frôles la transcendance."
		case g.BestClicks < 900000:
			g.phrase = "Le million est tout proche… peux-tu sentir le pouvoir te frôler ?"
		case g.BestClicks >= 1000000:
			g.phrase = "Tu as brisé la barrière du million… le monde te regarde. Es-tu prêt à renaître plus fort ?"
		default:
			g.phrase = "Tu n’es qu’à quelques clics du million… et d’une renaissance éternelle."
		}

		opchance := &text.DrawOptions{}
		opchance.GeoM.Translate(0, 10)
		opchance.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 215, B: 0, A: 255})
		text.Draw(screen, g.phrase, &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   12,
		}, opchance)
	} else {
		// 🌌 Fond doré et animé pour le Rebirth Store
		angle := float64((time.Now().UnixNano() / 5e6) % 360)
		for _, b := range g.BagsBG {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-32, -32)
			op.GeoM.Rotate(angle * math.Pi / 180.0)
			op.GeoM.Scale(b.scale, b.scale)
			op.GeoM.Translate(b.x, b.y)
			op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 80})
			screen.DrawImage(g.MoneyBag, op)
		}

		// 💎 Overlay doré léger
		ebitenutil.DrawRect(screen, 0, 0, 640, 480, color.RGBA{255, 215, 0, 25})

		// 💀 Panneau central du magasin
		ebitenutil.DrawRect(screen, 70, 40, 500, 400, color.RGBA{25, 25, 25, 220})
		ebitenutil.DrawRect(screen, 70, 40, 500, 5, color.RGBA{255, 215, 0, 255})
		ebitenutil.DrawRect(screen, 70, 435, 500, 5, color.RGBA{255, 215, 0, 255})

		// ✨ Titre stylé
		opTitle := &text.DrawOptions{}
		opTitle.GeoM.Translate(160, 70)
		opTitle.ColorScale.ScaleWithColor(color.RGBA{255, 215, 0, 255})
		text.Draw(screen, "⚡ Rebirth Soul Store ⚡", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   16,
		}, opTitle)

		// 💰 Icône et compteur d’âmes
		ebitenutil.DrawRect(screen, 90, 110, 460, 40, color.RGBA{60, 60, 60, 255})
		opSouls := &text.DrawOptions{}
		opSouls.GeoM.Translate(110, 135)
		opSouls.ColorScale.ScaleWithColor(color.RGBA{173, 216, 230, 255})
		text.Draw(screen, fmt.Sprintf("Âmes : %d", g.Âmes), &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   14,
		}, opSouls)

		// 💠 Offres visuelles (3 rangées hypnotiques)
		offers := []string{
			"Instinct automatique(Auto-clicks → +jauge)",
			"Renaissance dorée(Tes Rebirth te font garder 10% de ton argent)",
		}

		for i, txt := range offers {
			y := 180 + i*40
			ebitenutil.DrawRect(screen, 100, float64(y), 440, 30, color.RGBA{40, 40, 40, 255})
			ebitenutil.DrawRect(screen, 100, float64(y), 440, 2, color.RGBA{255, 215, 0, 100})
			opOffer := &text.DrawOptions{}
			opOffer.GeoM.Translate(120, float64(y+22))
			opOffer.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
			if i == 1 {
				text.Draw(screen, txt, &text.GoTextFace{
					Source: mplusFaceSource,
					Size:   6,
				}, opOffer)
			} else {
				text.Draw(screen, txt, &text.GoTextFace{
					Source: mplusFaceSource,
					Size:   10,
				}, opOffer)
			}
		}

		// 🔙 Bouton quitter
		ebitenutil.DrawRect(screen, 260, 400, 120, 30, color.RGBA{255, 215, 0, 150})
		opExit := &text.DrawOptions{}
		opExit.GeoM.Translate(280, 420)
		opExit.ColorScale.ScaleWithColor(color.RGBA{0, 0, 0, 255})
		text.Draw(screen, "Quitter", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   12,
		}, opExit)

	}
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
