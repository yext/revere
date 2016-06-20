# Data Sources

##### Path
"/datasources"

##### Requires
One or more data source(s)

##### Tests
**View**
1. All Data Sources
	* Load the page
	* Ensure that the graphite threshold data sources load with a valid URL field
**Edit**
1. Add a data source
	* Click Add
	* Type in a URL
	* Click Save
	* Reload the page - your new data source should still be there
2. Delete a data source
	* Click the "x" next to your data source
	* Your data source should disappear
	* Reload the page - your data source should not appear
