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

//nolint:funlen
func (s GameScore) CalculateTP() {
	type score struct {
		id int
		ps PlayerScore
	}

	scores := make([]score, 0, len(s))
	for id, ps := range s {
		scores = append(scores, score{id: id, ps: ps})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].ps.GP < scores[j].ps.GP
	})

	// Extract values once
	// gpN/idN - N means players place
	gp4, gp3, gp2, gp1 := scores[0].ps.GP, scores[1].ps.GP, scores[2].ps.GP, scores[3].ps.GP
	id4, id3, id2, id1 := scores[0].id, scores[1].id, scores[2].id, scores[3].id

	switch {
	// 4 < 3 < 2 < 1
	case gp4 < gp3 && gp3 < gp2 && gp2 < gp1:
		entry4 := s[id4]
		entry4.TP = fourthPlaceTP
		s[id4] = entry4
		entry3 := s[id3]
		entry3.TP = thirdPlaceTP
		s[id3] = entry3
		entry2 := s[id2]
		entry2.TP = secondPlaceTP
		s[id2] = entry2
		entry1 := s[id1]
		entry1.TP = firstPlaceTP
		s[id1] = entry1

	// 4 = 3 < 2 < 1
	case gp4 == gp3 && gp3 < gp2 && gp2 < gp1:
		tp43 := (fourthPlaceTP + thirdPlaceTP) / 2
		entry4 := s[id4]
		entry4.TP = tp43
		s[id4] = entry4
		entry3 := s[id3]
		entry3.TP = tp43
		s[id3] = entry3
		entry2 := s[id2]
		entry2.TP = secondPlaceTP
		s[id2] = entry2
		entry1 := s[id1]
		entry1.TP = firstPlaceTP
		s[id1] = entry1

	// 4 < 3 = 2 < 1
	case gp4 < gp3 && gp3 == gp2 && gp2 < gp1:
		tp32 := (thirdPlaceTP + secondPlaceTP) / 2
		entry4 := s[id4]
		entry4.TP = fourthPlaceTP
		s[id4] = entry4
		entry3 := s[id3]
		entry3.TP = tp32
		s[id3] = entry3
		entry2 := s[id2]
		entry2.TP = tp32
		s[id2] = entry2
		entry1 := s[id1]
		entry1.TP = firstPlaceTP
		s[id1] = entry1

	// 4 < 3 < 2 = 1
	case gp4 < gp3 && gp3 < gp2 && gp2 == gp1:
		tp21 := (secondPlaceTP + firstPlaceTP) / 2
		entry4 := s[id4]
		entry4.TP = fourthPlaceTP
		s[id4] = entry4
		entry3 := s[id3]
		entry3.TP = thirdPlaceTP
		s[id3] = entry3
		entry2 := s[id2]
		entry2.TP = tp21
		s[id2] = entry2
		entry1 := s[id1]
		entry1.TP = tp21
		s[id1] = entry1

	// 4 = 3 = 2 < 1
	case gp4 == gp3 && gp3 == gp2 && gp2 < gp1:
		tp432 := (fourthPlaceTP + thirdPlaceTP + secondPlaceTP) / 3
		entry4 := s[id4]
		entry4.TP = tp432
		s[id4] = entry4
		entry3 := s[id3]
		entry3.TP = tp432
		s[id3] = entry3
		entry2 := s[id2]
		entry2.TP = tp432
		s[id2] = entry2
		entry1 := s[id1]
		entry1.TP = firstPlaceTP
		s[id1] = entry1

	// 4 < 3 = 2 = 1
	case gp4 < gp3 && gp3 == gp2 && gp2 == gp1:
		tp321 := (thirdPlaceTP + secondPlaceTP + firstPlaceTP) / 3
		entry4 := s[id4]
		entry4.TP = fourthPlaceTP
		s[id4] = entry4
		entry3 := s[id3]
		entry3.TP = tp321
		s[id3] = entry3
		entry2 := s[id2]
		entry2.TP = tp321
		s[id2] = entry2
		entry1 := s[id1]
		entry1.TP = tp321
		s[id1] = entry1

	// 4 = 3 = 2 = 1
	case gp4 == gp3 && gp3 == gp2 && gp2 == gp1:
		tp4321 := (fourthPlaceTP + thirdPlaceTP + secondPlaceTP + firstPlaceTP) / 4
		for id := range s {
			entry := s[id]
			entry.TP = tp4321
			s[id] = entry
		}
	}
}
