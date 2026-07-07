package utils

func GenerateHandlerKey(method string, path string) string {
	return method + ":" + path
}
