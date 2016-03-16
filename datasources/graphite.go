package datasources

import (
	"encoding/json"
)

type GraphiteDataSource struct {
	Url string `json:"url"`
}

func init() {
	addDataSourceType(GraphiteDataSource{})
}

func (gs GraphiteDataSource) Id() DataSourceTypeId {
	return 0
}

func (gs GraphiteDataSource) Template() string {
	return "graphite-datasource.html"
}

func (gs GraphiteDataSource) Name() string {
	return "Graphite"
}

func (gs GraphiteDataSource) Scripts() []string {
	return []string{
		"graphite-datasource.js",
	}
}

func (gs GraphiteDataSource) LoadInfo(dataSourceInfo string) (parsedInfo interface{}, err error) {
	parsedInfo = new(GraphiteDataSource)
	err = json.Unmarshal([]byte(dataSourceInfo), &parsedInfo)
	return
}

func (gs GraphiteDataSource) DefaultInfo() interface{} {
	newSource := new(GraphiteDataSource)
	newSource.Url = ""
	return newSource
}
