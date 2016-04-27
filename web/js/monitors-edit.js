$(document).ready(function() {
  monitorsEdit.init();
});

var monitorsEdit = function() {
  var m = {};

  m.init = function() {
    triggersEdit.init();
    labelMonitors.init();
    initProbe();
    initForm();
  };

  var initProbe = function() {
    var $probeType = $('#js-probe-type');
    $probeType.change(function() {
      var $error = $('.js-error').first().empty();
      $error.addClass('hidden');
      $.ajax({
        url: '/monitors/new/probe/edit/' + encodeURIComponent($probeType.val()),
        contentType: 'application/json; charset=UTF-8'
      }).success(function(response) {
        if (response.template) {
          $('#js-probe').html(response.template);
        }
      }).fail(function(jqXHR, textStatus, errorThrown) {
        $('#js-errors').html($error);
        $error.append(jqXHR.responseText).removeClass('hidden');
      });
    });
  };

  var initForm = function() {
    $('#js-monitor-form').submit(function(e) {
      e.preventDefault();
      var $form = $(this);
      var url = $form.attr('action'),
        data = $.extend(
          getMonitorData(),
          {'ProbeParams': JSON.stringify(getProbeData())},
          {'Triggers': triggersEdit.getData()},
          {'Labels': labelMonitors.getData()}
        );
      $.ajax({
        url: url,
        method: 'POST',
        data: JSON.stringify(data),
        contentType: 'application/json; charset=UTF-8'
      }).success(function(response) {
        if (response.errors) {
          return revere.showErrors(response.errors);
        }
        if (response.redirect) {
          window.location.replace(response.redirect);
        } else {
          window.location.replace('/monitors/' + data['id']);
        }
      }).fail(function(jqXHR, textStatus, errorThrown) {
        revere.showErrors([jqXHR.responseText || textStatus]);
      });
    });
  }

  // Serializing functions
  var getMonitorData = function() {
    return $('#js-monitor-info').find(':input').serializeObject();
  };

  var getProbeData = function() {
    return $('#js-probe').find(':input').serializeObject();
  };

  return m;
}();
