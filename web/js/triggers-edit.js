var triggersEdit = function() {
  var tse = {};

  tse.init = function() {
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
      var id = $trigger.find('input[name="TriggerID"]').val();
      if(id == '0'){
        $trigger.remove();
      } else {
        $trigger.addClass('hidden');
        $trigger.find('input[name="Delete"]').prop('checked', true);
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

  return tse;
}();
