$(document).ready(function() {
  datasources.addSourceFunction(graphiteDataSourceHandler.getData);
});


var graphiteDataSourceHandler = function() {
  var gdsh = {}

  gdsh.getData = function() {
    var data = [];
    $.each($('.js-datasource.graphite'), function() {
      var sendData = $(this).find(':input.required').serializeObject();
      var sourceData = $(this).find(':input.source').serializeObject();
      $.extend(sendData, {'DataSourceParams': JSON.stringify(sourceData)});
      data.push(sendData)
    });
    return data;
  };

  return gdsh
}();
