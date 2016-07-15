$(document).ready(function() {
  graphiteThreshold.init();
});

var graphiteThreshold = function() {
  var g = {};

  g.init = function() {
    addSerializeFn();
  };

  var addSerializeFn = function() {
    probes.addSerializeFn($('#js-graphite-threshold-probe-type').val(), function(probe) {
      var inputs = probe.find(':input:not(.js-threshold)').serializeObject();
      probe.find(':input.js-threshold').each(function() {
          if ($(this).val() == "") {
              $(this).remove();
          }
      });
      thresholds = probe.find(':input.js-threshold').serializeObject();

      return JSON.stringify($.extend(inputs, {"Thresholds": thresholds}));
    });
  };

  return g;
}();
