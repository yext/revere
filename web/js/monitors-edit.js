$(document).ready(function() {
  monitorsEdit.init();
});

var probes = function() {
  var p = {};
  var fns = {};

  p.addSerializeFn = function(probeType, fn) {
    fns[probeType] = fn;
  };

  p.getSerializeFn = function(probeType) {
    return fns[probeType];
  };

  return p;
}();

var monitorsEdit = function() {
  var m = {};

  m.init = function() {
    monitorTriggersEdit.init();
    monitorLabels.init();
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
          {'Triggers': monitorTriggersEdit.getData()},
          {'Labels': monitorLabels.getData()}
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
    var monitorInputs = $('#js-monitor-info').find(':input').serializeObject(),
      probeFn = probes.getSerializeFn(monitorInputs['ProbeType']),
      probe;
    if (probeFn !== undefined) {
      probe = probeFn($('#js-probe'));
    }

    if (probe === undefined) {
      probe = JSON.stringify($('#js-probe').find(':input').serializeObject());
    }
    return $.extend(monitorInputs, {'ProbeParams':probe});

  };

  return m;
}();
