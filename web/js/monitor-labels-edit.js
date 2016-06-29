var monitorLabelsEdit = function() {
  var ml = {},
    cl = componentList('monitor-label');

  ml.init = function() {
    cl.init();
  };

  ml.getData = function() {
    var clData = cl.getData();
      data = [];
    $.each(clData, function(i, monitorLabel) {
      monitorLabel.Label = {
        'LabelID': monitorLabel.LabelID
      };
      delete monitorLabel.LabelID

      monitorID = parseInt($('input[name=MonitorID]').val())
      data.push($.extend(monitorLabel, {'MonitorID': monitorID}));
    });
    return data;
  };

  return ml;
}();
