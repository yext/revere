var labelMonitorsEdit = function() {
  var lm = {},
    cl = componentList('label-monitor');

  lm.init = function() {
    cl.init();
  };

  lm.getData = function() {
    var clData = cl.getData();
      data = [];
    $.each(clData, function(i, labelMonitor) {
      labelMonitor.Monitor = {
        'MonitorID': labelMonitor.MonitorID
      }
      delete labelMonitor.MonitorID

      labelID = parseInt($('input[name=LabelID]').val())
      data.push($.extend(labelMonitor, {'LabelID': labelID}));
    });
    return data;
  };

  return lm;
}();
