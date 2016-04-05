$(document).ready(function() {
  subprobesIndex.init();
});

var subprobesIndex = function() {
  var s = {};

  s.init = function() {
    enteredStates.init();
  };

  return s;
}();
