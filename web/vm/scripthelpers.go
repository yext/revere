package vm

import "path"

const (
	baseDir         = "web/js"
	baseServingPath = "static/js"
)

func GetScript(filepath string) string {
	return path.Join(baseServingPath, filepath)
}

func AppendDir(dir string, scripts []string) []string {
	result := make([]string, len(scripts))
	for i, script := range scripts {
		result[i] = path.Join(dir, script)
	}
	return result
}
