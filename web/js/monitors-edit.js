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
      url: '/monitors/new/probe/' + encodeURIComponent($probeType.val()),
      contentType: 'application/json; charset=UTF-8',
    }).success(function(response) {
      if (response.template) {
        $('#probe').html(response.template);
      }
    }).fail(function(jqXHR, textStatus, errorThrown) {
      // 500
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
      contentType: 'application/json; charset=UTF-8',
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
      // 500
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
  var data = []
  $.each($('.trigger'), function() {
    var triggerInputs = $(this).find(':input').serializeObject();
    data.push(triggerInputs);
  });
  return data;
}
