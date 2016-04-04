$(document).ready(function() {
  activeIssues.init();
});

var activeIssues = function() {
  var ai = {};

  ai.init = function() {
    labelFilter.initForm();
  };
  return ai;
}();
