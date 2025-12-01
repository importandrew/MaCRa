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
