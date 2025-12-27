package models

type Wind int

const (
	East Wind = iota
	South
	West
	North
)

func (w Wind) Kanji() string {
	windToKanji := map[Wind]string{
		East:  "東",
		South: "南",
		West:  "西",
		North: "北",
	}
	return windToKanji[w]
}

var Winds = []Wind{East, South, West, North}

var RoundWinds = map[int][]Wind{
	// East
	1: {East, South, West, North},
	2: {South, West, North, East},
	3: {West, North, East, South},
	4: {North, East, South, West},

	// South
	5: {South, East, North, West},
	6: {East, North, West, South},
	7: {North, West, South, East},
	8: {West, South, East, North},

	// West
	9:  {West, North, South, East},
	10: {North, South, East, West},
	11: {South, East, West, North},
	12: {East, West, North, South},

	// North
	13: {North, West, East, South},
	14: {West, East, South, North},
	15: {East, South, North, West},
	16: {South, North, West, East},
}

var RoundSeats = map[int][]Wind{
	// East
	1: {East, South, West, North},
	2: {East, South, West, North},
	3: {East, South, West, North},
	4: {East, South, West, North},

	// South
	5: {South, East, North, West},
	6: {South, East, North, West},
	7: {South, East, North, West},
	8: {South, East, North, West},

	// West
	9:  {West, North, South, East},
	10: {West, North, South, East},
	11: {West, North, South, East},
	12: {West, North, South, East},

	// North
	13: {North, West, East, South},
	14: {North, West, East, South},
	15: {North, West, East, South},
	16: {North, West, East, South},
}
