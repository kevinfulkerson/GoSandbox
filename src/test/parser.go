package main

// parseLookupString parses the provided string for an id element to use for creating a lookup array.
// It returns the first set of ordered values as a slice, or the entire string if no lookup value was found.
// TODO: Incorporate this searching into a parsing utility that will read lines and parse out each element (CONT.)
// as it is encountered, lookup the correct value and then insert that value as a string into the original
// parsed contents. This should also allow for looking up local values simply using a @id:local/.. value,
// and external values using @id:/.. value.
func ParseLookupString(lookupString string) (lookupValues []string) {
	tempString := ""
	elementStarted := false
	for i := 0; i < len(lookupString); i++ {
		switch lookupString[i] {
		case '@':
			// In most cases, this will be the start of an element id indication, so
			// look for the rest of the indicator. The indicator is currently:
			// @id:
			addToString := true
			if i+3 <= len(lookupString) {
				if lookupString[i+1] == 'i' && lookupString[i+2] == 'd' && lookupString[i+3] == ':' {
					addToString = false
					elementStarted = true
					i += 3
				}
			}

			if addToString {
				tempString += string(lookupString[i])
			}
		case '/':
			if elementStarted {
				// Check if we have anything to insert, and if so, insert it.
				if tempString != "" {
					lookupValues = append(lookupValues, tempString)
					tempString = ""
				}
			} else {
				tempString += string(lookupString[i])
			}
		case '[':
			if elementStarted {
				// This is an array indication, so we need to insert the index
				// after its parent. To do this, insert the parent first.
				lookupValues = append(lookupValues, tempString)

				// Set the temp value to the empty string so the next '/' token
				// we encounter (after the ending array indicator) won't cause an
				// insert of a blank value in the array.
				tempString = ""

				// The parent is in correctly, so begin inserting the index value.
				str := ""
				i++
				for lookupString[i] != ']' {
					str += string(lookupString[i])
					i++
				}

				lookupValues = append(lookupValues, str)
			} else {
				tempString += string(lookupString[i])
			}
		default:
			if lookupString[i] != '\n' && lookupString[i] != '\r' {
				tempString += string(lookupString[i])
			}
		}
	}

	return append(lookupValues, tempString)
}
