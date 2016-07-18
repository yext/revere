$(document).ready(function() {
  graphiteDSLoader.init();
});

var graphiteDSLoader = function() {
  var gdl = {};

  gdl.init = function() {
    var probeType = $('#js-probe-type option:selected').val();
    $.ajax({
      url: '/datasources/probe/' + probeType,
    }).done(function(data, status, jqXHR) {
      gdl.displayDataSources(JSON.parse(data));
    }).fail(function(jqXHR, status, error) {
      revere.showErrors([error]);
    });
  };

  gdl.displayDataSources = function(datasources) {
    var $selector = $('#js-datasources'),
      selectedUrl = $selector.data('url');
    $.each(datasources, function(i, datasource) {
      var url = datasource.DataSource.URL,
        selected = url === selectedUrl,
        id = datasource.SourceID;

      $selector.append($('<option></option').html(url)
        .data('id',id).attr('selected', selected));
    });
  };

  return gdl;
}();
