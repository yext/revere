package web

import (
	"errors"
)

// Provides a way to pass multiple values to a subtemplate
// Taken from: http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters
func dict(args ...interface{}) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("Invalid dict call")
	}
	dict := make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			return nil, errors.New("Dict keys must be strings")
		}
		dict[key] = args[i+1]
	}
	return dict, nil
}

func isLastBc(a []breadcrumb, i int) bool {
	return i == len(a)-1
}
