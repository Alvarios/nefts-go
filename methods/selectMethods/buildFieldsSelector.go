package selectMethods

import "strings"

// User can restrict query to some specific fields. Ignored fields will not be sent in the response request.
func buildFieldsSelector(fields []string, selector string) string {
	fieldsString := ""

	if fields != nil && len(fields) > 0 {
		for _, field := range fields {
			// Add separators between fields declarations, as required by N1QL syntax.
			if len(fieldsString) > 0 {
				fieldsString += ", "
			}

			// 'all' selects all NON META fields.
			if field == "all" || field == "*" {
				fieldsString += selector + ".*"
			} else if strings.HasPrefix(field, "meta.") && len(field) > 5 {
				fieldsString += "META(" + selector + ")." + field[5:]
			} else {
				fieldsString += field
			}
		}
	} else {
		// Default selection.
		fieldsString = selector + ".*, META(" + selector + ").id"
	}

	return fieldsString
}
