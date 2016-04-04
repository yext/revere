var labelFilter = function() {
  var lf = {};
  lf.initForm = function() {
    var labelId = revere.getParameterByName("label");
    if (labelId === null) return;
    $('.labels').each(function() {
      if ($(this).val() === labelId) {
        $(this).prop('selected', true);
      }
    });
  };
  return lf;
}();
