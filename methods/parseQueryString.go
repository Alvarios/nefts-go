package methods

import "strings"

func parser(query string, labels map[string][]string) map[string]string {
	// quoteFlag : tells parser we are currently inside quotes, so we ignore default break rules.
	// ignoreCurrentChar : don't insert current char in parsed string.
	// index : last group index.
	// mapped : a list of captured groups of words.
	quoteFlag, ignoreCurrentChar, index, mapped := false, false, 0, []string{""}

	// Loop through each character. This method is more complex but allows more stable behavior.
	for i, char := range query {
		// Get last mapped group of characters.
		last := mapped[index]

		// Ignore breaks inside quotations.
		if char == ' ' && !quoteFlag {
			ignoreCurrentChar = true

			// If a previous group was captured, add a new one. Empty previous groups can occur when a label is set
			// without any content.
			if len(last) > 0 {
				mapped = append(mapped, "")
				index++
			}
		}

		// Open a quote.
		if char == '"' && !quoteFlag {
			// Only open if a closing quote is present.
			for _, sc := range query[i + 1:] {
				if sc == '"' {
					quoteFlag = true
					ignoreCurrentChar = true
					break
				}
			}
		// Close a quote.
		} else if char == '"' && quoteFlag {
			quoteFlag = false
			ignoreCurrentChar = true
		}

		// Add char to current group.
		if ignoreCurrentChar == false {
			mapped[index] += string(char)
		}

		ignoreCurrentChar = false
	}

	// Build the output object that groups groups by label.
	output := map[string]string{"general": ""}
	// key : current label key to assign the group to.
	// matchFlag : check, if a label-like indication is present, if it matches any actually available label.
	// cutLen : position of the ":" label mark in the current group, if any.
	key, matchFlag, cutLen := "general", false, 0

	// Loop through groups of characters previously captured. We now want to group and assign them to a label.
	for _, str := range mapped {
		// We assume no label match prior to any operation.
		matchFlag = false

		// Loop through labels to find if any match the current label mark.
		for fieldKey, value := range labels {
			// A field can be associated with multiple label marks.
			for _, prefix := range value {
				// A matching label mark was found.
				if strings.HasPrefix(str, prefix + ":") {
					matchFlag = true
					cutLen = len(prefix) + 1
					break
				}
			}

			// Update current key value.
			if matchFlag {
				key = fieldKey
				break
			} else {
				key = "general"
				cutLen = 0
			}
		}

		// Check if current key is already present in output. If so, add the current value to the previously computed
		// one.
		if val, ok := output[key]; ok {
			output[key] = val + " " + str[cutLen:]
		} else {
			output[key] = str[cutLen:]
		}
	}

	return output
}

func ParseQueryString(query string, labels map[string][]string) map[string]string {
	if len(query) > 0 {
		if labels != nil && len(labels) > 0 {
			// Split string according to labels.
			return parser(query, labels)
		} else {
			// No labels to parse.
			return map[string]string{"general": query}
		}
	}

	return nil
}
