package main

func joinList(list []string) string {
	joinedString := ""
	for _, item := range list {
		joinedString += item + "\n"
	}

	return joinedString
}
