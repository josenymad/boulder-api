package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/josenymad/boulder-api/config"
)

func GetNumberOfRounds(competition string) (count int, err error) {
	query := "SELECT COUNT(*) FROM rounds WHERE competition_id = (SELECT competition_id FROM competitions WHERE competition_name = $1)"
	err = config.DB.QueryRow(query, competition).Scan(&count)
	if err != nil {
		return 0, errors.New("failed to get round count")
	}
	return count, nil
}

func RemoveLastComma(query string) string {
	lastIndex := strings.LastIndex(query, ",")

	if lastIndex != -1 {
		query = query[:lastIndex] + query[lastIndex+1:]
	}

	return query
}

func BuildScoresQueryString(category string, competition string) (query string, err error) {
	queryStart := "SELECT competitor_name, SUM(points) AS TOTAL,\n"
	queryEnd := fmt.Sprintf(
		`FROM (
		SELECT 
			c.name AS competitor_name,
			cat.name AS category_name,
			r.round_number,
			s.points
		FROM 
			competitors c
		INNER JOIN 
			competition_categories cat ON c.category_id = cat.category_id
		LEFT JOIN 
			scores s ON c.competitor_id = s.competitor_id
		LEFT JOIN 
			boulder_problems bp ON s.problem_id = bp.problem_id
		LEFT JOIN 
			rounds r ON bp.round_id = r.round_id
		WHERE
			r.competition_id = (SELECT competition_id FROM competitions WHERE competition_name = '%s')
	) AS subquery
	WHERE
		category_name = '%s'
	GROUP BY
		competitor_name
	ORDER BY
		TOTAL DESC;`,
		competition,
		category,
	)

	numberOfRounds, err := GetNumberOfRounds(competition)
	if err != nil {
		return query, errors.New("failed to get number of rounds")
	}

	var iteratedQuery string
	for index := 1; index <= numberOfRounds; index++ {
		iteratedQuery += fmt.Sprintf("SUM(CASE WHEN round_number = %d THEN points ELSE 0 END) AS round_%d,\n", index, index)
	}

	iteratedQuery = RemoveLastComma(iteratedQuery)

	query = queryStart + iteratedQuery + queryEnd

	return query, err
}
