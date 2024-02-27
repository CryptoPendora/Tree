func PrintDebug(format string, a ...interface{}) {
	if !UnDebugTextMode {
		fmt.Printf(format, a...)
	}
}
