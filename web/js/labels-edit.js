$(document).ready(function() {
  labelsEdit.init();
});

var labelsEdit = function() {
  var le = {};

  le.init = function() {
    labelTriggersEdit.init();
    labelMonitorsEdit.init();
    initForm();
  };

  var initForm = function() {
    $('#js-label-form').submit(function(e) {
      e.preventDefault();
      var $form = $(this);
      var url = $form.attr('action');

      data = $.extend(
        getLabelData(),
        {'Monitors': labelMonitorsEdit.getData()},
        {'Triggers': labelTriggersEdit.getData()}
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
          window.location.replace('/labels/' + data['LabelID']);
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
