package resource

import (
	"encoding/json"

	"github.com/yext/revere/db"
)

type Graphite struct{}

type GraphiteResource struct {
	Graphite
	URL string
}

// Eventually implemented in DB layer
type GraphiteResourceDBModel struct {
	URL string
}

func init() {
	addType(Graphite{})
}

func (Graphite) Id() db.ResourceType {
	return 0
}

func (Graphite) Name() string {
	return "Graphite"
}

func (Graphite) loadFromParams(ds string) (Resource, error) {
	var g GraphiteResource
	err := json.Unmarshal([]byte(ds), &g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (Graphite) loadFromDB(ds string) (Resource, error) {
	var g GraphiteResourceDBModel
	err := json.Unmarshal([]byte(ds), &g)
	if err != nil {
		return nil, err
	}

	return &GraphiteResource{
		URL: g.URL,
	}, nil
}

func (Graphite) blank() (Resource, error) {
	return &GraphiteResource{}, nil
}

func (Graphite) Templates() string {
	return "graphite-resource.html"
}

func (Graphite) Scripts() []string {
	return []string{
		"graphite-resource.js",
	}
}

func (g GraphiteResource) Serialize() (string, error) {
	gDB := GraphiteResourceDBModel{
		g.URL,
	}

	gDBJSON, err := json.Marshal(gDB)
	return string(gDBJSON), err
}

func (g GraphiteResource) Type() ResourceType {
	return Graphite{}
}

func (g GraphiteResource) Validate() []string {
	var errs []string
	if g.URL == "" {
		errs = append(errs, "Url is required")
	}

	return errs
}
