package game

func (g *Game) Players() [4]int {
	return g.startSeating
}

func (g *Game) RoundStartPlayers() [4]int {
	return g.roundSeating
}

func (g *Game) Opposite(id int) int {
	return g.currentSeating[(g.currentSeating.Wind(id)+2)%4]
}

func (g *Game) Right(id int) int {
	return g.currentSeating[(g.currentSeating.Wind(id)+1)%4]
}

func (g *Game) Left(id int) int {
	return g.currentSeating[(g.currentSeating.Wind(id)+3)%4]
}
