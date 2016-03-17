$(document).ready(function() {
  labelsEdit.init();
});

var labelsEdit = function() {
  var le = {};

  le.init = function() {
    initLabelTriggers();
    initLabelMonitors();
    initForm();
  };

  var initLabelTriggers = function() {
    triggersEdit.init();
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

    var $baseAddMonitor = $('.js-new-label-monitor').first();

    $.each($('.js-label-monitor'), function() {
      var id = $(this).find('input[name="id"]').val();
      hideOption(id, $('.js-monitor-name'));
    });

    $('#js-add-monitor').click(function(e) {
      e.preventDefault();
      var $existingMonitors = $('.js-label-monitor'),
        $newMonitors = $('.js-new-label-monitor'),
        $newMonitor = $baseAddMonitor.clone();

      if ($newMonitors.length > 1) {
        $newMonitor.insertAfter($newMonitors.last());
      } else {
        $newMonitor.insertAfter($existingMonitors.last());
      }
      $newMonitor.removeClass('hidden');
    });

    $(document.body).on('click', '.js-remove-label-monitor', function(e) {
      e.preventDefault();
      $monitor = $(this).parents('.js-label-monitor');
      $monitor.hide();
      $monitor.find('input[name="delete"]').prop('checked', true);
      showOption($monitor.find('input[name="id"]').val(), $('.js-monitor-name'));
    });

    $(document.body).on('click', '.js-remove-new-label-monitor', function(e) {
      e.preventDefault();
      $this = $(this),
        $newMonitor = $this.parents('.js-new-label-monitor'),
        $selectedOption = $newMonitor.find('.js-monitor-name option:selected');

      showOption($selectedOption.val(), $('.js-monitor-name'));
      $newMonitor.remove();
    });

    $(document.body).on('change', '.js-monitor-name', function() {
      var $this = $(this),
        $selectedOption = $this.find('option:selected'),
        $otherSelectBoxes = $('.js-monitor-name').not($this);

      var previousId = $this.data('id'),
        description = $selectedOption.data('description');

      $this.parents('.js-new-label-monitor').find('.description').text(description);
      $this.data('id', $selectedOption.val());

      showOption(previousId, $otherSelectBoxes);
      hideOption($selectedOption.val(), $otherSelectBoxes);
    });
  };

  var initForm = function() {
    $('#js-label-form').submit(function(e) {
      e.preventDefault();
      var $form = $(this);
      var url = $form.attr('action');

      data = $.extend(
        getLabelData(),
        {'monitors': getLabelMonitorData()},
        {'triggers': triggersEdit.getData()}
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
          window.location.replace('/labels/' + data['id']);
        }
      }).fail(function(jqXHR, textStatus, errorThrown) {
        var $error = $('.js-error').first().empty();
        $('#js-errors').html($error);
        $error.append(jqXHR.responseText).removeClass('hidden');
      });
    });
  };


  var getLabelData = function() {
    return $('#js-label-info').find(':input').serializeObject();
  };

  var getLabelMonitorData = function() {
    var data = [];
    $.each($.merge($('.js-label-monitor'),
        $('.js-new-label-monitor:visible')), function() {
      data.push($(this).find(':input').serializeObject());
    });
    return data;
  };

  return le;
}();
