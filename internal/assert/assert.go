package assert

// Assert expected states and world views..
func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}
