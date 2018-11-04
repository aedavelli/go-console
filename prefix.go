package console

func Prefix() (string, bool) {
	const separator = " > "
	p := appName + separator
	if presentCtx != "" {
		p += presentCtx + separator
	}
	return p, true
}
