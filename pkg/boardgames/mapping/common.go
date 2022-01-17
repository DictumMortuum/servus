package mapping

func Hamming(s1 string, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	var column []bool

	if len(r1) >= len(r2) {
		column = make([]bool, len(r1)+1)

		for i := 0; i < len(r2); i++ {
			column[i] = r1[i] == r2[i]
		}
	} else {
		column = make([]bool, len(r2)+1)

		for i := 0; i < len(r1); i++ {
			column[i] = r1[i] == r2[i]
		}
	}

	distance := 0

	for _, item := range column {
		if item == true {
			distance += 1
		}
	}

	return distance
}
