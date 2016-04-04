$(document).ready(function() {
  monitorsIndex.init();
});

var monitorsIndex = function() {
  var m = {};

  m.init = function() {
    labelFilter.initForm();
  };
  return m;
}();
