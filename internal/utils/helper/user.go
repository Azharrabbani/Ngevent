package helper

func IsAuthorized(id, compareID string) bool {
	if id != compareID {
		return false
	}

	return true
}
