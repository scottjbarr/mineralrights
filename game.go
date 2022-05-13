package mineralrights

import (
	"fmt"

	"math/rand"
	"strings"
	"time"
)

type Game struct {
	Mines        int64   // L
	Workers      int64   // P
	Money        int64   // M
	FoodPrice    int64   // FP
	OrePerMine   int64   // CE
	Storage      int64   // C
	Satisfaction float64 // S
	Year         int64   // Y
	MinePrice    int64   // LP
	OrePrice     int64   // CP
}

func New() *Game {
	rand.Seed(time.Now().Unix())

	g := Game{
		Mines:        int64(rnd(1)*3 + 5),
		Workers:      int64(rnd(1)*66 + 40),
		FoodPrice:    int64(rnd(1)*40 + 80),
		OrePerMine:   int64(rnd(1)*40 + 80),
		Storage:      0,
		Satisfaction: 1.0,
		Year:         1,
		MinePrice:    getMineMarketPrice(),
		OrePrice:     getOreMarketPrice(),
	}

	g.Money = int64(rnd(1)*50+10) * g.Workers

	return &g
}

func Run() {
	title()

	for {
		g := New()

		loop(g)

		// play again?
		if !playAgain() {
			break
		}
	}
}

func loop(g *Game) {
	for {
		g.SummaryScreen()

		//
		// inputs from the player
		//
		g.oreSell()

		g.minesSell()

		if g.Mines == 0 {
			g.gameOver("NO MORE MINES....GAME OVER.")
			return
		}

		fmt.Printf("\nBUYING\n")
		fmt.Printf("------\n\n")

		// how much food to buy?
		purchased := g.foodBuy()

		g.satisfactionAdjust(purchased)

		// how many mines to buy?
		g.minesBuy()

		if !g.satisfactionCheck() {
			g.gameOver("THE WORKERS REVOLTED!")
			return
		}

		// events and disasters
		//
		// check that there are enough workers
		if !g.checkWorkers() {
			g.gameOver("YOU'VE OVERWORKED EVERYONE!")
			return
		}

		// has a radiation leak occurred?
		g.radiationLeak()

		// do we have a market glut?
		g.marketGlut()

		// increment the year
		g.Year += 1

		// set buy/sell prices
		g.MinePrice = getMineMarketPrice()
		g.OrePrice = getOreMarketPrice()
	}
}

func (g *Game) radiationLeak() {
	if g.Workers < 30 {
		g.adjustWorkers()
		return
	}

	if rnd(1) > 0.01 {
		return
	}

	fmt.Printf("RADIOACTIVE LEAK...........MANY DIE!\n")

	g.Workers /= 2
}

func (g *Game) marketGlut() {
	if g.OrePerMine < 150 {
		return
	}

	fmt.Printf("MARKET GLUT..........PRICE DROPS!\n")

	g.OrePerMine /= 2
}

func (g *Game) checkWorkers() bool {
	if g.Workers/g.Mines < 10 {
		return false
	}

	g.adjustWorkers()

	return true
}

func (g *Game) adjustWorkers() {
	if g.Satisfaction > 1.1 {
		g.Workers += int64(rnd(1)*10 + 1)
	}

	if g.Satisfaction < 0.9 {
		g.Workers -= int64(rnd(1)*10 + 1)
	}
}

func (g *Game) satisfactionAdjust(foodPurchased int64) {
	if foodPurchased/g.Workers > 120 {
		g.Satisfaction += 0.1
	}

	if foodPurchased/g.Workers < 80 {
		g.Satisfaction -= 0.1
	}
}

func (g *Game) satisfactionCheck() bool {
	if g.Satisfaction < 0.6 {
		return false
	}

	if g.Satisfaction > 1.1 {
		g.OrePerMine += int64(rnd(1)*20 + 1)
	}

	if g.Satisfaction < 0.9 {
		g.OrePerMine -= int64(rnd(1)*20 + 1)
	}

	return true
}

func (g *Game) oreSell() {
	v := buildOreSellValidator(g.Storage)

	i := readIntWithPrompt("HOW MUCH ORE DO YOU WISH TO SELL?", v)

	// TAKES AWAY SOLD ORE
	g.Storage += i

	// ADDS TO MONEY SUPPLY
	g.Money += i * g.OrePrice

	if i > 0 {
		g.printBalance()
	}
}

func (g *Game) minesSell() {
	v := buildMinesSellValidator(g.Mines)

	i := readIntWithPrompt("HOW MANY MINES DO YOU WISH TO SELL?", v)

	// TAKES AWAY MINE(S)
	g.Mines -= i

	// ADDS TO MONEY SUPPLY
	g.Money += i * g.MinePrice

	if i > 0 {
		g.printBalance()
	}
}

func (g *Game) minesBuy() {
	v := func(i int64) bool {
		return i >= 0
	}

	i := readIntWithPrompt("HOW MANY MINES DO YOU WISH TO BUY?", v)

	// INCREASE NO. OF MINES IF NEEDED
	g.Mines += i

	// ADJUST MONEY SUPPLY AGAIN
	g.Money -= i * g.MinePrice

	if i > 0 {
		g.printBalance()
	}
}

func (g *Game) foodBuy() int64 {
	v := func(i int64) bool {
		return i >= 0 && i <= g.Money
	}

	i := readIntWithPrompt("HOW MUCH TO SPEND ON FOOD (APPR.$100 EA.)?", v)

	// ADJUSTS MONEY SUPPLY
	g.Money -= i

	if i > 0 {
		g.printBalance()
	}

	return i
}

func (g *Game) printBalance() {
	fmt.Printf("YOU HAVE $%0v\n", g.Money)
}

func (g *Game) SummaryScreen() {
	// current state affairs of colony
	fmt.Printf("\nYEAR %v\n\n", g.Year)

	fmt.Printf("THERE ARE %v WORKERS IN THE COLONY.\n", g.Workers)
	fmt.Printf("YOU HAVE %v MINES AND $%v\n", g.Mines, g.Money)
	fmt.Printf("SATISFACTION FACTOR IS %0.1f\n", g.Satisfaction)

	fmt.Printf("YOUR MINES PRODUCED %v TONS EACH.\n", g.OrePerMine)

	g.Storage += g.OrePerMine * g.Mines

	fmt.Printf("AMOUNT OF ORE IN STORE IS %v TONS\n\n", g.Storage)

	fmt.Printf("SELLING\n")
	fmt.Printf("-------\n\n")
	fmt.Printf("ORE SELLING PRICE IS $%v PER TON.\n", g.OrePrice)
	fmt.Printf("MINE SELLING PRICE IS $%v PER MINE.\n", g.MinePrice)
}

func (g *Game) gameOver(message string) {
	fmt.Printf("%s\n", message)

	fmt.Printf("\nGAME OVER\n\n")
	fmt.Printf("\n\nYOU LASTED %v YEARS\n", g.Year)
}

func playAgain() bool {
	fmt.Printf("\n\n\nWOULD YOU LIKE TO HAVE ANOTHER GAME?\n")
	fmt.Printf("\n(Y/N) ")

	v := readString()

	// anything except "y" is considered a signal to end
	return strings.ToLower(strings.TrimSpace(v)) == "y"
}

func getMineMarketPrice() int64 {
	return int64(rnd(1)*2000 + 2000)
}

func getOreMarketPrice() int64 {
	return int64(rnd(1)*12 + 7)
}

func title() {
	fmt.Printf("\n\nMINERAL RIGHTS\n")
	fmt.Printf("--------------\n\n")

	fmt.Printf("(c) 1985 GN Woodhead, Wyke, Bradford.\n")
	fmt.Printf("Published in \"Your Computer\" magazine July 1985\n")
	fmt.Printf("Ported by Scott Barr <scottjbarr@gmail.com> 2022.\n\n")

	fmt.Printf("YOU ARE THE NEWLY ELECTED LEADER OF A\n")
	fmt.Printf("MINING COLONY ON THE PLANET ASTRON.\n\n")
	fmt.Printf("ALL DECISIONS CONCERNING THE SALE\n")
	fmt.Printf("OF ORE TO INTERGALACTIC TRADES, FOOD\n")
	fmt.Printf("PURCHASES AND BUYING & SELLING OF\n")
	fmt.Printf("MINES ARE MADE BY YOU!\n\n")
	fmt.Printf("THERE MUST BE AT LEAST 10 WORKERS/MINE.\n\n")
}
