package wire

func generateHandlerKey(method string, path string) string {
	return method + ":" + path
}
