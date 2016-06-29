# Monitors

##### Path:
"/Monitors/<id>"
"/Monitors/<id>/edit"

##### Requires:
One or more Monitor(s)
Optional: One or more Label(s)
Optional: One or more Trigger(s)

##### Tests:
**View:**

1. Single Monitor
  * Load the page
  * Ensure that the Monitor fields and associated Labels and Triggers fields appear

**Edit:**

1. Edit Monitor attribute
  * Load the page
  * Modify an attribute of the Monitor (e.g. Name, Description)
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
2. Edit Probe attribute
  * Load the page
  * Modify an attribute of the Probe (e.g. Thresholds, Graphite Data Source)
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
3. Create Trigger
  * Load the page
  * Click "+ Add" under Triggers
  * Fill out appropriate fields for the new Trigger
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
4. Delete Trigger
  * Load the page
  * Click the "x" next to the newly created Trigger
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
5. Add Label
  * Load the page
  * Click "+ Add" under Labels
  * Fill out appropriate fields
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
6. Remove Label
  * Load the page
  * Click the "x" next to the newly added Label
  * Click "Save"
  * Changes should be reflected on the view page after the redirect  
  _Note: the Label itself should still exist, but will not be associated with the Monitor_
7. Validation
  * Load the page
  * Enter in an invalid input for a field (e.g. empty string for Monitor name)
  * Click "Save"
  * Error should pop up in red at the top of the screen
