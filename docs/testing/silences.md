# Silences

##### Path:
"/silences/<id>"

##### Requires:
One Monitor
One Silence applied to the Monitor

##### Tests:
**View:**

1. Single Silence
  * Load the page
  * Ensure the Monitor name, subprobes, start, and end are displayed

**Edit**

1. Create Silence
  * Go to ```/silences/new``` or click the "+new" button on the index page
  * Select Monitor and duration of Silence
  * Click Save
  * Go to index page, new Silence should show up
2. Edit Silence - Datetime Picker
  * Create or go to a future Silence
  * Edit the start using the datetime picker
  * Edit the end using the datetime picker
  * Click Save
  * Changes should be reflected in the view Silence page after the redirect
3. Edit Silence - Now
  * Create or go to a future Silence
  * Edit the start by clicking the "Now" radio button
  * Click Save
  * Changes should be reflected in the view Silence page after the redirect
4. Edit Silence - End Now
  * Create or go to a current Silence
  * Click the "End Now" button
  * Silence end should be set to current time in the view Silence page, and should appear as a past Silence in the index
