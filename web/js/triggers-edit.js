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

var triggersEdit = function() {
  var te = {};

  te.init = function() {
    initTriggers();
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

  te.getData = function() {
    var data = []
    $.each($('.js-trigger'), function() {
      var triggerOptions = $(this).find('.js-trigger-options :input').serializeObject(),
        targetFn = targets.getSerializeFn(triggerOptions['targetType']),
        target;
      if (targetFn !== undefined) {
        target = targetFn($(this).find('.js-target'));
      }

      if (target === undefined) {
        target = JSON.stringify($(this).find('.js-target :input').serializeObject());
      }
      data.push($.extend(triggerOptions, {'target':target}));
    });
    return data;
  };

  return te;
}();
