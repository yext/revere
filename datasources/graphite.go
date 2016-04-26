package datasources

import (
	"encoding/json"
)

type Graphite struct{}

type GraphiteDataSource struct {
	Url string
}

func init() {
	addDataSourceType(Graphite{})
}

func (gs Graphite) Id() DataSourceTypeId {
	return 0
}

func (gs Graphite) Template() string {
	return "graphite-datasource.html"
}

func (gs Graphite) Name() string {
	return "Graphite"
}

func (gs Graphite) Scripts() []string {
	return []string{
		"graphite-datasource.js",
	}
}

func (gs Graphite) Load(dataSourceJson string) (dataSource DataSource, err error) {
	dataSource = new(GraphiteDataSource)
	err = json.Unmarshal([]byte(dataSourceJson), &dataSource)
	return
}

func (gs Graphite) LoadDefault() DataSource {
	newSource := new(GraphiteDataSource)
	newSource.Url = ""
	return newSource
}

func (g *GraphiteDataSource) Validate() []string {
	var errs []string
	if g.Url == "" {
		errs = append(errs, "Url is required")
	}

	return errs
}

func (g *GraphiteDataSource) DataSourceType() DataSourceType {
	return Graphite{}
}
