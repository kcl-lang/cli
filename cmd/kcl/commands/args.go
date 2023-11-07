package cmd

func argsGet(a []string, n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}
