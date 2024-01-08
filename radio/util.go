package radio

func nullTerminatedBytesToString(b []byte) string {
	// Find the index of the null terminator
	nullIndex := -1
	for i, v := range b {
		if v == 0 {
			nullIndex = i
			break
		}
	}

	// If null terminator is not found, return the whole slice as string
	if nullIndex == -1 {
		return string(b)
	}

	// Convert the relevant portion to string
	str := string(b[:nullIndex])

	return str
}
