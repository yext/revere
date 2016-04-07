# revere
“One if by land, and two if by sea”—Alerting for Graphite

## Configuration

### Adding Data Source Types

To add a data source type, you'll need to add three files:

1. A go file in datasources/, which implements DataSourceType (in datasources/datasource)
2. A go template in web/views/datasources/
3. At least one js file in web/js/datasources

#### Go file

This file should implement DataSourceType - see graphite.go for an example. the `Template`
function should return the name of the go template file you create, and the Scripts function should
return an array containing the names of the js files you create.

The LoadInfo function should parse the JSON string stored in the string passed to it into a struct representation that can be accessed by a custom Probe.

#### Go template

This file will determine how data sources of your new type are displayed/edited on the data sources page in Revere. Within this file, you will have access to the view model described in web/vm/datasourcetype.go. General notes:

* Each data source should have its own row
* Each data source should be enclosed in a div with the `js-datasource` class
* Each data source should have a delete button, with the class `delete`
* This template needs to have an add button, which should have a class `js-add-source`
* The add button needs to have a `data-sourceref` attribute unique to this data source type - in addition, each element that has `js-datasource` as a class should also have this attribute value as a class.

See web/views/datasources/graphite-datasource.html for an example.

#### JS Files

In this file you should create a module with a function that serializes all data sources of your type into a json with format `{id: Number, sourceTypeId: Number, source: String}`. Then, call `datasources.addSourceFunction(<your function here>)` in a `$(document).ready()` callback. See web/js/datasources/graphite-datasource.js for an example.
