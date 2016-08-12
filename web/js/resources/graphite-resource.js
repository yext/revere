$(document).ready(function() {
  resources.addSourceFunction(graphiteResourceHandler.getData);
});


var graphiteResourceHandler = function() {
  var gdsh = {}

  gdsh.getData = function() {
    var data = [];
    $.each($('.js-resource.graphite'), function() {
      var sendData = $(this).find(':input.required').serializeObject();
      var sourceData = $(this).find(':input.source').serializeObject();
      $.extend(sendData, {'ResourceParams': JSON.stringify(sourceData)});
      data.push(sendData)
    });
    return data;
  };

  return gdsh
}();
