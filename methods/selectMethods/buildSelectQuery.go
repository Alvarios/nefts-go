package selectMethods

import "github.com/Alvarios/nefts-go/config"

func BuildSelectQuery(
	query map[string]string,
	bucket string,
	fields []string,
	joins []config.Join,
	labelOptions map[string]config.LabelOption,
) (string, *config.Error) {
	output := "SELECT "

	// Selector will identify the main bucket in the request.
	selector := "b"

	// If joins are present, we reserve the 'b' selector for later.
	areJoinsPresent := joins != nil && len(joins) > 0
	if areJoinsPresent {
		selector = "`" + bucket + "`"
	}

	fieldsSelector := buildFieldsSelector(fields, selector)
	scoreSelector := buildScoreSelector(query, labelOptions)
	joinsSelector, err := buildJoinsFieldsSelector(joins)

	if err != nil {
		return "", err
	}

	links, err := linkJoins(joins, selector)

	if err != nil {
		return "", err
	}

	// Joins query has a different structure from basic query.
	if joins != nil && len(joins) > 0 {
		output = "SELECT b.*"
		if len(scoreSelector) > 0 {
			output += ", " + scoreSelector
		}

		output += " \n\tFROM(\n\t\tSELECT " + fieldsSelector
		if len(joinsSelector) > 0 {
			output += ", " + joinsSelector
		}

		output += " FROM `" + bucket + "`" + links + "\n\t) AS b"
	} else {
		output += fieldsSelector
		if len(scoreSelector) > 0 {
			output += ", " + scoreSelector
		}

		output += "\n\tFROM `" + bucket + "` AS b"
	}

	return output, nil
}
