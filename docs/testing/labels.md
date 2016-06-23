# Labels

##### Path:
View: ```/labels/<id>```
Edit: ```/labels/<id>/edit```

##### Requires:
One or more Label(s)
One or more Trigger(s) on the Label
One or more Monitor(s) associated with Label

##### Tests:
**View:**

1. Labels
  * Load the page
  * Make sure the triggers and monitors appear properly

**Edit:**

1. Edit Label attribute
  * Load the page
  * Modify an attribute of the Label (e.g. Name, Description)
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
2. Create Label Trigger
  * Load the page
  * Click "+ Add" under Triggers
  * Fill out appropriate fields for the new Trigger
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
3. Delete Label Trigger
  * Load the page
  * Click the "x" next to the newly created Trigger
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
4. Add Monitor
  * Load the page
  * Click "+ Add" under Monitors
  * Fill out appropriate fields
  * Click "Save"
  * Changes should be reflected on the view page after the redirect
5. Remove Monitor
  * Load the page
  * Click the "x" next to the newly added Monitor
  * Click "Save"
  * Changes should be reflected on the view page after the redirect  
  _Note: the Monitor itself should still exist, but will not be associated with the label_
