var targets = function() {
  var t = {};
  var fns = {};

  t.addSerializeFn = function(targetType, fn) {
    fns[targetType] = fn;
  };

  t.getSerializeFn = function(targetType) {
    return fns[targetType];
  };

  return t;
}();

$(document).ready(function() {
  var $triggers = $('.trigger');
  // Don't allow the removal of triggers if there is only one
  if ($triggers.length == 1) {
    $('.remove-trigger').hide();
  }

  /*
   * Event handlers
   */
  $('#add-trigger').click(function(e) {
    e.preventDefault();
    $('.remove-trigger').show();
    $triggers.first()
      .clone()
      .insertAfter('.trigger:last')
      .find('input[type="text"]')
      .val('');
  });

  $(document.body).on('click', '.remove-trigger', function(e) {
    e.preventDefault();
    if ($('.trigger').length > 1) {
      $(this).parents('.trigger').remove();
    }
    if ($('.trigger').length == 1) {
      $('.remove-trigger').hide();
    }
  });

  var $probeType = $('#probe-type');
  $probeType.change(function() {
    var $error = $('.error').first().empty();
    $error.addClass('hidden');
    $.ajax({
      url: '/monitors/new/probe/edit/' + encodeURIComponent($probeType.val()),
      contentType: 'application/json; charset=UTF-8'
    }).success(function(response) {
      if (response.template) {
        $('#probe').html(response.template);
      }
    }).fail(function(jqXHR, textStatus, errorThrown) {
      $('#errors').html($error);
      $error.append(jqXHR.responseText).removeClass('hidden');
    });
  });

  $('.targetType').change(function() {
    var $error = $('.error').first().empty();
    $error.addClass('hidden');
    $that = $(this);
    $.ajax({
      url: '/monitors/new/target/edit/' + encodeURIComponent($that.val()),
      contentType: 'application/json; charset=UTF-8'
    }).success(function(response) {
      if (response.template) {
        $that.parents('.form-group')
          .next('.target')
          .html(response.template);
      }
    }).fail(function(jqXHR, textStatus, errorThrown) {
      $('#errors').html($error);
      $error.append(jqXHR.responseText).removeClass('hidden');
    });
  });

  $('#monitor-form').submit(function(e) {
    e.preventDefault();
    var $form = $(this);

    var url = $form.attr('action'),
      data = $.extend(
        getMonitorData(),
        {'probe': JSON.stringify(getProbeData())},
        {'triggers': getTriggerData()});

    $.ajax({
      url: url,
      method: 'POST',
      data: JSON.stringify(data),
      contentType: 'application/json; charset=UTF-8'
    }).success(function(response) {
      if (response.errors) {
        var $error = $('.error').first().empty();
        $('#errors').html($error);
        $.each(response.errors, function() {
          $error.append(this + '<br/>').removeClass('hidden');
        });
        return;
      }
      if (response.redirect) {
        window.location.replace(response.redirect);
      } else {
        window.location.replace('/monitors/' + data['id']);
      }
    }).fail(function(jqXHR, textStatus, errorThrown) {
      var $error = $('.error').first().empty();
      $('#errors').html($error);
      $error.append(jqXHR.responseText).removeClass('hidden');
    });
  });
});

var getMonitorData = function() {
  return $('#monitor-info').find(':input').serializeObject();
}

var getProbeData = function() {
  return $('#probe').find(':input').serializeObject();
}

var getTriggerData = function() {
  var data = [];
  $.each($('.trigger'), function() {
    var triggerOptions = $(this).find('.trigger-options :input').serializeObject(),
      targetFn = targets.getSerializeFn(triggerOptions['targetType']),
      target;
    if (targetFn !== undefined) {
      target = targetFn($(this).find('.target'));
    }

    if (target === undefined) {
      target = $(this).find('.target :input').serializeObject();
    }
    data.push($.extend(triggerOptions, {'target':target}));
  });
  return data;
}
