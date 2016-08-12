# revere
“One if by land, and two if by sea”—Alerting for Graphite

## Configuration

### Adding Resource Types

To add a resource type, you'll need to add three files:

1. A go file in resources/, which implements ResourceType (in resources/resource)
2. A go template in web/views/resources/
3. At least one js file in web/js/resources

#### Go file

This file should implement ResourceType - see graphite.go for an example. the `Template`
function should return the name of the go template file you create, and the Scripts function should
return an array containing the names of the js files you create.

The LoadDefault function should parse the JSON string stored in the string passed to it into a struct representation that can be accessed by a custom Probe.

#### Go template

This file will determine how resources of your new type are displayed/edited on the resources page in Revere. Within this file, you will have access to the view model described in web/vm/resourcetype.go. General notes:

* Each resource should have its own row
* Each resource should be enclosed in a div with the `js-resource` class, and the name of the resource type
* Each resource should have a delete button, with the class `delete`, and a checkbox

See web/views/resources/graphite-resource.html for an example.

#### JS Files

In this file you should create a module with a function that serializes all resources of your type into a json with format `{resourceId: Number, resourceTypeId: Number, resource: String}`. Then, call `resources.addSourceFunction(<your function here>)` in a `$(document).ready()` callback. See web/js/resources/graphite-resource.js for an example.
