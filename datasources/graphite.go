package datasources

import "encoding/json"

type Graphite struct{}

// TODO(fchen): match front-end to "URL"
type GraphiteDataSource struct {
	URL string
}

// Eventually implemented in DB layer
type GraphiteDataSourceDBModel struct {
	URL string
}

func init() {
	addDataSourceType(Graphite{})
}

func (_ Graphite) Id() DataSourceTypeId {
	return 0
}

func (_ Graphite) Name() string {
	return "Graphite"
}

func (_ Graphite) loadFromParams(ds string) (DataSource, error) {
	var g GraphiteDataSource
	err := json.Unmarshal([]byte(ds), &g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (_ Graphite) loadFromDb(ds string) (DataSource, error) {
	var g GraphiteDataSourceDBModel
	err := json.Unmarshal([]byte(ds), &g)
	if err != nil {
		return nil, err
	}

	return &GraphiteDataSource{
		g.URL,
	}, nil
}

func (_ Graphite) blank() (DataSource, error) {
	return &GraphiteDataSource{}, nil
}

func (_ Graphite) Templates() string {
	return "graphite-datasource.html"
}

func (_ Graphite) Scripts() []string {
	return []string{
		"graphite-datasource.js",
	}
}

func (g GraphiteDataSource) Serialize() (string, error) {
	gDB := GraphiteDataSourceDBModel{
		g.URL,
	}

	gDBJSON, err := json.Marshal(gDB)
	return string(gDBJSON), err
}

// TODO(fchen): check for and fix references to DataSourceType in frontend
func (g GraphiteDataSource) Type() DataSourceType {
	return Graphite{}
}

func (g GraphiteDataSource) Validate() []string {
	var errs []string
	if g.URL == "" {
		errs = append(errs, "Url is required")
	}

	return errs
}
