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

  var noTriggers = function() {
    return $('.js-trigger').not('.hidden').length == 0;
  }

  var initTriggers = function() {
    var $triggers = $('.js-trigger'),
      $emptyTriggers = $('.js-empty-triggers'),
      $triggerTemplate = $triggers.first();

    $.each($('.js-trigger').not(':first'), function() {
      $(this).removeClass('hidden');
    });

    if(noTriggers()) {
      $emptyTriggers.removeClass('hidden');
    }

    $('#js-add-trigger').click(function(e) {
      e.preventDefault();
      $newTrigger = $triggerTemplate
        .clone()
        .insertAfter('.js-trigger:last');
      $newTrigger.removeClass('hidden');
      $emptyTriggers.addClass('hidden');
    });

    $(document.body).on('click', '.js-remove-trigger', function(e) {
      e.preventDefault();
      $trigger = $(this).parents('.js-trigger');
      var id = $trigger.find('input[name="id"]').val();
      if(id == '0'){
        $trigger.remove();
      } else {
        $trigger.addClass('hidden');
        $trigger.find('input[name="delete"]').prop('checked', true);
      }
      if(noTriggers()) {
        $emptyTriggers.removeClass('hidden');
      }
    });

    $('.js-targetType').change(function() {
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
        revere.showErrors([jqXHR.responseText || textStatus]);
      });
    });
  };

  te.getData = function() {
    var data = []
    $.each($('.js-trigger').not(':first'), function() {
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
