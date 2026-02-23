package usecase

import (
	"context"
	"sort"

	"xyz-football-api/internal/modules/report"

	"github.com/google/uuid"
)

func preComputeTeamWins(matches []report.MatchRow) map[uuid.UUID]int {
	wins := make(map[uuid.UUID]int)
	for _, m := range matches {
		if m.HomeScore > m.AwayScore {
			wins[m.HomeTeamID]++
		} else if m.AwayScore > m.HomeScore {
			wins[m.AwayTeamID]++
		}
	}
	return wins
}

type topScorerPerMatch struct {
	name string
}

func preComputeTopScorers(goals []report.GoalEvent) map[uuid.UUID]topScorerPerMatch {
	type playerInfo struct {
		name  string
		count int
	}
	matchPlayers := make(map[uuid.UUID]map[uuid.UUID]*playerInfo)

	for _, g := range goals {
		if _, ok := matchPlayers[g.MatchID]; !ok {
			matchPlayers[g.MatchID] = make(map[uuid.UUID]*playerInfo)
		}
		pi, exists := matchPlayers[g.MatchID][g.PlayerID]
		if !exists {
			pi = &playerInfo{name: g.PlayerName}
			matchPlayers[g.MatchID][g.PlayerID] = pi
		}
		pi.count++
	}

	result := make(map[uuid.UUID]topScorerPerMatch, len(matchPlayers))
	for matchID, players := range matchPlayers {
		var topName string
		topGoals := 0
		for _, pi := range players {
			if pi.count > topGoals {
				topGoals = pi.count
				topName = pi.name
			}
		}
		result[matchID] = topScorerPerMatch{name: topName}
	}

	return result
}

func determineFinalStatus(homeScore, awayScore int) string {
	if homeScore > awayScore {
		return report.FinalStatusHomeWin
	} else if homeScore < awayScore {
		return report.FinalStatusAwayWin
	}
	return report.FinalStatusDraw
}

func (u *reportUsecase) GetMatchReports(ctx context.Context, page, limit int) ([]report.MatchReportResponse, int, error) {
	matches, err := u.reportRepo.GetFinishedMatches(ctx)
	if err != nil {
		return nil, 0, err
	}

	goals, err := u.reportRepo.GetGoalEvents(ctx)
	if err != nil {
		return nil, 0, err
	}

	teamWins := preComputeTeamWins(matches)
	topScorers := preComputeTopScorers(goals)

	results := make([]report.MatchReportResponse, 0, len(matches))

	for _, m := range matches {
		topScorer := "-"
		if ts, ok := topScorers[m.ID]; ok {
			topScorer = ts.name
		}

		r := report.MatchReportResponse{
			MatchID:           m.ID.String(),
			MatchDate:         m.MatchDate.Format("2006-01-02"),
			MatchTime:         m.MatchTime,
			HomeTeam:          m.HomeTeamName,
			AwayTeam:          m.AwayTeamName,
			HomeScore:         m.HomeScore,
			AwayScore:         m.AwayScore,
			FinalStatus:       determineFinalStatus(m.HomeScore, m.AwayScore),
			TopScorer:         topScorer,
			HomeTeamTotalWins: teamWins[m.HomeTeamID],
			AwayTeamTotalWins: teamWins[m.AwayTeamID],
		}
		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].MatchDate == results[j].MatchDate {
			return results[i].MatchTime > results[j].MatchTime
		}
		return results[i].MatchDate > results[j].MatchDate
	})

	totalData := len(results)
	start := (page - 1) * limit
	if start >= totalData {
		return []report.MatchReportResponse{}, totalData, nil
	}
	end := min(start+limit, totalData)
	// end := start + limit
	// if end > totalData {
	// 	end = totalData
	// }

	return results[start:end], totalData, nil
}
