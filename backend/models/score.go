package models

import (
	"bytes"
	"sort"
	"text/template"
)

const (
	firstPlaceTP  float32 = 4
	secondPlaceTP float32 = 2
	thirdPlaceTP  float32 = 1
	fourthPlaceTP float32 = 0
)

type GameScore map[int]PlayerScore

type PlayerScore struct {
	GP int     // game points
	TP float32 // tournament points
}

func (s GameScore) String() string {
	var output bytes.Buffer
	_ = template.Must(template.New("score").Parse(
		`{{range $player, $score := .}}{{$player}}={{$score.GP}}{{"\t"}}{{end}}`,
	)).Execute(&output, s)
	return output.String()
}

func (s GameScore) Final() string {
	var output bytes.Buffer
	_ = template.Must(template.New("final").Parse(
		`{{range $player, $score := .}}{{$player}}={{$score.GP}}({{$score.TP}}){{"\t"}}{{end}}`,
	)).Execute(&output, s)
	return output.String()
}

// CalculateTP sets each player's tournament points from their place by game points.
// Players with equal game points share the average TP of the places they occupy.
func (s GameScore) CalculateTP() {
	placeTP := [...]float32{firstPlaceTP, secondPlaceTP, thirdPlaceTP, fourthPlaceTP}

	ids := make([]int, 0, len(s))
	for id := range s {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return s[ids[i]].GP > s[ids[j]].GP
	})

	for start := 0; start < len(ids); {
		end := start + 1
		for end < len(ids) && s[ids[end]].GP == s[ids[start]].GP {
			end++
		}

		var sum float32
		for place := start; place < end && place < len(placeTP); place++ {
			sum += placeTP[place]
		}
		tp := sum / float32(end-start)

		for _, id := range ids[start:end] {
			ps := s[id]
			ps.TP = tp
			s[id] = ps
		}
		start = end
	}
}
