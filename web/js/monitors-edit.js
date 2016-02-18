$(document).ready(function() {
  monitorsEdit.init();
});

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

var monitorsEdit = function() {
  var m = {};

  m.init = function() {
    initProbe();
    initTriggers();
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

  var initTriggers = function() {
    var $triggers = $('.js-trigger');

    $('#js-add-trigger').click(function(e) {
      e.preventDefault();
      $newTrigger = $triggers.first()
        .clone()
        .insertAfter('.js-trigger:last')
      $newTrigger.find('input[type="text"]').val('');
      $newTrigger.find('input[name="id"]').val(0);
      $newTrigger.find('input[name="delete"]').prop('checked', false);
      $newTrigger.show();
    });

    $(document.body).on('click', '.js-remove-trigger', function(e) {
      e.preventDefault();
      $trigger = $(this).parents('.js-trigger');
      $trigger.hide();
      $trigger.find('input[name="delete"]').prop('checked', true);
    });

    $('.js-targetType').change(function() {
      var $error = $('.js-error').first().empty();
      $error.addClass('hidden');
      $that = $(this);
      $.ajax({
        url: '/monitors/new/target/edit/' + encodeURIComponent($that.val()),
        contentType: 'application/json; charset=UTF-8'
      }).success(function(response) {
        if (response.template) {
          $that.parents('.js-trigger-options')
            .next('.js-target')
            .html(response.template);
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
          {'probe': JSON.stringify(getProbeData())},
          {'triggers': getTriggerData()}
        );
      $.ajax({
        url: url,
        method: 'POST',
        data: JSON.stringify(data),
        contentType: 'application/json; charset=UTF-8'
      }).success(function(response) {
        if (response.errors) {
          var $error = $('.js-error').first().empty();
          $('#js-errors').html($error);
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
        var $error = $('.js-error').first().empty();
        $('#js-errors').html($error);
        $error.append(jqXHR.responseText).removeClass('hidden');
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

  var getTriggerData = function() {
    var data = []
    $.each($('.js-trigger'), function() {
      var triggerOptions = $(this).find('.js-trigger-options :input').serializeObject(),
        targetFn = targets.getSerializeFn(triggerOptions['targetType']),
        target;
      if (targetFn !== undefined) {
        target = targetFn($(this).find('.js-target'));
      }

      if (target === undefined) {
        target = $(this).find('.js-target :input').serializeObject();
      }
      data.push($.extend(triggerOptions, {'target':target}));
    });
    return data;
  };

  return m;
}();
