$(document).ready(function() {
  monitorsIndex.init();
});

var monitorsIndex = function() {
  var m = {};

  m.init = function() {
    initFilterForm();
  };

  var initFilterForm = function() {
    var labelId = revere.getParameterByName("label");
    if (labelId === null) return;
    $('.labels').each(function() {
      if ($(this).val() === labelId) {
        $(this).prop('selected', true);
      }
    });
  };
  return m;
}()
