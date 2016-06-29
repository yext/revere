var monitorTriggersEdit = function() {
  var mte = {};

  mte.init = function() {
    triggersEdit.init();
  };

  mte.getData = function() {
    var data = [];
    $.each($('.js-trigger').not(':first'), function() {
      var monitorTrigger = {
        Trigger: triggerEdit.getData(this),
        MonitorID: parseInt($('input[name="MonitorID"]').val())
      };
      monitorTrigger.Subprobes = monitorTrigger.Trigger.Subprobes;
      delete monitorTrigger.Trigger.Subprobes;
      data.push(monitorTrigger);
    });
    return data;
  };

  return mte;
}();
