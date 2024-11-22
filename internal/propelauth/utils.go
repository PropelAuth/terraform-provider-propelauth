package propelauth

func Contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func FlipBoolRef(b *bool) *bool {
	if b == nil {
		return nil
	} else if *b {
		new_b := false
		return &new_b
	} else {
		new_b := true
		return &new_b
	}
}
