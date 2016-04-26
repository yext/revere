$(document).ready(function() {
  labelsEdit.init();
});

var labelsEdit = function() {
  var le = {};

  le.init = function() {
    initLabelTriggers();
    labelMonitors.init();
    initForm();
  };

  var initLabelTriggers = function() {
    triggersEdit.init();
  };

  var initForm = function() {
    $('#js-label-form').submit(function(e) {
      e.preventDefault();
      var $form = $(this);
      var url = $form.attr('action');

      data = $.extend(
        getLabelData(),
        {'monitors': labelMonitors.getData()},
        {'Triggers': triggersEdit.getData()}
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
          window.location.replace('/labels/' + data['id']);
        }
      }).fail(function(jqXHR, textStatus, errorThrown) {
        revere.showErrors([jqXHR.responseText || textStatus]);
      });
    });
  };


  var getLabelData = function() {
    return $('#js-label-info').find(':input').serializeObject();
  };

  return le;
}();
