$(document).ready(function() {
  graphiteResourceLoader.init();
});

var graphiteResourceLoader = function() {
  var gdl = {};

  gdl.init = function() {
    var probeType = $('#js-probe-type option:selected').val();
    $.ajax({
      url: '/resources/probe/' + probeType,
    }).done(function(data, status, jqXHR) {
      gdl.displayResources(JSON.parse(data));
    }).fail(function(jqXHR, status, error) {
      revere.showErrors([error]);
    });
  };

  gdl.displayResources = function(resources) {
    var $selector = $('#js-resources'),
      selectedUrl = $selector.data('url');
    $.each(resources, function(i, resource) {
      var url = resource.Resource.URL,
        selected = url === selectedUrl,
        id = resource.ResourceID;

      $selector.append($('<option></option').html(url)
        .data('id',id).attr('selected', selected));
    });
  };

  return gdl;
}();
