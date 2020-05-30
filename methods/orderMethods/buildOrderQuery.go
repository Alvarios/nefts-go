package orderMethods

import (
	"github.com/Alvarios/nefts-go/config"
	"github.com/Alvarios/nefts-go/utils"
	"strings"
)

func BuildOrderQuery(
	query map[string]string,
	order map[string]string,
	labelOptions map[string]config.LabelOption,
) string {
	output := ""

	if query != nil && len(query) > 0 {
		output += "ORDER BY ("

		subOutput := ""

		for key, value := range query {
			if len(utils.SoftFormat(value)) > 0 {
				if len(subOutput) > 0 {
					subOutput += " + "
				}

				labelParams := config.LabelOption{}

				// Search if some custom parameter was provided for current label
				if labelOptions != nil {
					match, ok := labelOptions[key]
					if ok {
						labelParams = match
					} else {
						match, ok = labelOptions["global"]
						if ok {
							labelParams = match
						}
					}
				}

				if labelParams.Weight == "" {
					labelParams.Weight = "1"
				}

				subOutput += utils.HardFormat(key) + "_score * " + labelParams.Weight
			}
		}

		output += subOutput + ") DESC"
	}

	if order != nil && len(order) > 0 {
		for key, value := range order {
			if len(output) > 0 {
				output += ", "
			} else {
				output += "ORDER BY "
			}

			output += "b." + key + " " + strings.ToUpper(value)
		}
	}

	return output
}
