var labelMonitors = function() {
  var lm = {};

  lm.init = function() {
    initLabelMonitors();
  };

  lm.getData = function() {
    var data = [];
    $.each($.merge($('.js-label-monitor'),
        $('.js-new-label-monitor:visible')), function() {
      data.push($(this).find(':input').serializeObject());
    });
    return data;
  };

  var initLabelMonitors = function() {
    var showOption = function(id, $selectBoxes) {
      $.each($selectBoxes, function() {
        $(this).find('option[value="' + id + '"]').removeClass('hidden');
      });
    };

    var hideOption = function(id, $selectBoxes) {
      if (id === "") {
        return;
      }
      $.each($selectBoxes, function() {
        $(this).find('option[value="' + id + '"]').addClass('hidden');
      });
    };

    var $baseAddLabelMonitor = $('.js-new-label-monitor').first();

    $.each($('.js-label-monitor'), function() {
      var id = $(this).find('input.js-id').val();
      hideOption(id, $('.js-label-monitor-name'));
    });

    $('#js-add-label-monitor').click(function(e) {
      e.preventDefault();
      var $existingMonitors = $('.js-label-monitor'),
        $newLabelMonitors = $('.js-new-label-monitor'),
        $newLabelMonitor = $baseAddLabelMonitor.clone();

      $('.js-empty-label-monitor').addClass('hidden');
      if ($newLabelMonitors.length > 1 || $existingMonitors.length === 0) {
        $newLabelMonitor.insertAfter($newLabelMonitors.last());
      } else {
        $newLabelMonitor.insertAfter($existingMonitors.last());
      }
      $newLabelMonitor.removeClass('hidden');
    });

    $(document.body).on('click', '.js-remove-label-monitor', function(e) {
      e.preventDefault();
      $monitor = $(this).parents('.js-label-monitor');
      $monitor.hide();
      $monitor.find('input[name="Delete"]').prop('checked', true);
      showOption($monitor.find('input.js-id').val(), $('.js-label-monitor-name'));
    });

    $(document.body).on('click', '.js-remove-new-label-monitor', function(e) {
      e.preventDefault();
      $this = $(this),
        $newLabelMonitor = $this.parents('.js-new-label-monitor'),
        $selectedOption = $newLabelMonitor.find('.js-label-monitor-name option:selected');

      showOption($selectedOption.val(), $('.js-label-monitor-name'));
      $newLabelMonitor.remove();
    });

    $(document.body).on('change', '.js-label-monitor-name', function() {
      var $this = $(this),
        $selectedOption = $this.find('option:selected'),
        $otherSelectBoxes = $('.js-label-monitor-name').not($this);

      var previousId = $this.data('id'),
        description = $selectedOption.data('description');

      $this.parents('.js-new-label-monitor').find('.description').text(description);
      $this.data('id', $selectedOption.val());

      showOption(previousId, $otherSelectBoxes);
      hideOption($selectedOption.val(), $otherSelectBoxes);
    });
  };

  return lm;
}();
