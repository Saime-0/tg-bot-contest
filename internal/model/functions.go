package model

func MemberIDs(m []Member) []int {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}

	return ids
}
