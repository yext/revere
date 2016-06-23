var labelTriggersEdit = function() {
  var lte = {};

  lte.init = function() {
    triggersEdit.init();
  };

  lte.getData = function() {
    var data = [];
    $.each($('.js-trigger').not(':first'), function() {
      var labelTrigger = {
        Trigger: triggerEdit.getData(this),
        LabelID: parseInt($('input[name=LabelID]').val())
      };
      data.push(labelTrigger);
    });
    return data;
  };

  return lte;
}();
