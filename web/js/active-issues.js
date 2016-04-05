$(document).ready(function() {
  activeIssues.init();
});

var activeIssues = function() {
  var ai = {};

  ai.init = function() {
    enteredStates.init();
    labelFilter.initForm();
  };
  return ai;
}();
