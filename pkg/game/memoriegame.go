package game

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"test2/pkg/arr"
	randomcat "test2/pkg/randomCat"
)

type MemoryGame struct {
	Players             map[string]*Player // player list
	Cards               []*Card            // cards list
	PlayerPlayDirection []string           // the name player in random order play turn
	ActivePlayerIndex   int                // player who is this turn
	CancelStartChanel   chan struct{}      `json:"-"` // toggle cooldown to start
}

type Card struct {
	Id     int    // a random cart id
	TwinId int    // and this random twin cardid
	Img    string // the img link
	IsFind bool   // if the both twin card are find
}

type CardByID []*Card

func (a CardByID) Len() int           { return len(a) }
func (a CardByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CardByID) Less(i, j int) bool { return a[i].Id < a[j].Id }

type Player struct {
	Id    int  // this uuid
	Score int  // this score
	Ready bool // ready to start game
}

func NewGame(name string) *MemoryGame { // create a intence of game with a set of card
	g := &MemoryGame{Players: make(map[string]*Player), CancelStartChanel: make(chan struct{})}
	go func() { g.Cards = CreatXTwinCard(15) }()
	return g
}

func CreatXTwinCard(nb int) (cards []*Card) {
	for i, c := range randomcat.GetXRandomCatImg(nb) {
		cards = append(cards, &Card{Id: i, Img: c})
		cards = append(cards, &Card{Id: i, Img: c})
	}
	idused := " , "
	rmax, rmin := nb*2, 1
	getrandom := func() int {
		if rmax == rmin {
			return rmax
		}
		for r := rand.Intn(rmax-rmin) + rmin; ; r = rand.Intn(rmax-rmin) + rmin {
			if rs := fmt.Sprint(r); !strings.Contains(idused, " , "+rs+" , ") {
				idused += rs + " , "
			recheck:
				if r == rmin || strings.Contains(idused, fmt.Sprint(" , ", rmin, " , ")) {
					rmin++
					goto recheck
				} else if r == rmax || strings.Contains(idused, fmt.Sprint(" , ", rmax, " , ")) {
					rmax--
					goto recheck
				}
				return r
			}
		}
	}
	for i := 0; i < len(cards); i += 2 {
		c1 := getrandom()
		c2 := getrandom()
		cards[i].Id, cards[i+1].Id = c1, c2
		cards[i].TwinId, cards[i+1].TwinId = c2, c1
	}
	c := CardByID(cards)
	sort.Sort(c)
	cards = []*Card(c)
	return
}

func (g *MemoryGame) SetPlayerOrderRandom() {
	g.PlayerPlayDirection = make([]string, 0) // len(g.Players))
	for name := range g.Players {
		g.PlayerPlayDirection = append(g.PlayerPlayDirection, name)
	}
	fmt.Println(g.PlayerPlayDirection)
}

// add player if not exist
func (g *MemoryGame) AddPlayer(pname string) {
	if _, ok := g.Players[pname]; !ok {
		g.Players[pname] = &Player{Id: len(g.Players)}
	}
}

func (g *MemoryGame) RemovePlayer(pname string) {
	delete(g.Players, pname)
	for i, pn := range g.PlayerPlayDirection {
		if pn == pname {
			g.PlayerPlayDirection = arr.RemoveIndex(g.PlayerPlayDirection, i)
			return
		}
	}
}

func (g *MemoryGame) GetPlayerByName(n string) *Player {
	return g.Players[n]
}

func (p *Player) ToggleReady() {
	p.Ready = !p.Ready
}

func (g *MemoryGame) IsAllPlayersReady() bool {
	for _, v := range g.Players {
		if !v.Ready {
			return false
		}
	}
	return true
}

func (g *MemoryGame) NextPlayerTurn() {
	if g.ActivePlayerIndex+1 < len(g.PlayerPlayDirection) {
		g.ActivePlayerIndex++
	} else {
		g.ActivePlayerIndex = 0
	}
}

func (g *MemoryGame) GetPlayerTurnName() string {
	return g.PlayerPlayDirection[g.ActivePlayerIndex]
}

func (g *MemoryGame) CheckTwin(c1, c2 int) bool {
	return g.Cards[c1-1].TwinId == c2
}
