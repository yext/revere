/*
Package boxes is a wrapper around a set of go.rice boxes created by a go.rice
Config defined to first look for files in the package directory, and then in
the binary.
*/
package boxes

import rice "github.com/GeertJohan/go.rice"

//go:generate rice embed-go

// There are some odd things in this package due to how go.rice's rice tool
// works.  In order to generate the appropriate go files, the tool must go
// through this package and find all calls to FindBox (and MustFindBox) to
// figure out which directories need to be put into "Boxes".

// go.rice idiosyncrasy #1 - The package name variable must shadow the Config
// (if a rice.Config is used) in order for the rice tool to pick up calls to
// FindBox.

// go.rice idiosyncrasy #2 - Calls to FindBox must be done with string
// literals.

var (
	favicon   *rice.Box
	templates *rice.Box
	css       *rice.Box
	js        *rice.Box

	revereConfig = rice.Config{
		LocateOrder: []rice.LocateMethod{
			rice.LocateFS,
			rice.LocateWorkingDirectory,
			rice.LocateEmbedded,
			rice.LocateAppended,
		},
	}
)

func Favicon() *rice.Box {
	if favicon != nil {
		return favicon
	}

	rice := revereConfig

	var err error
	favicon, err = rice.FindBox("../web/favicon")
	if err != nil {
		panic(err)
	}

	return favicon
}

func Templates() *rice.Box {
	if templates != nil {
		return templates
	}

	rice := revereConfig

	var err error
	templates, err = rice.FindBox("../web/views")
	if err != nil {
		panic(err)
	}

	return templates
}

func CSS() *rice.Box {
	if css != nil {
		return css
	}

	rice := revereConfig

	var err error
	css, err = rice.FindBox("../web/css")
	if err != nil {
		panic(err)
	}

	return css
}

func JS() *rice.Box {
	if js != nil {
		return js
	}

	rice := revereConfig

	var err error
	js, err = rice.FindBox("../web/js")
	if err != nil {
		panic(err)
	}

	return js
}
