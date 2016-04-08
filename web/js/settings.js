$(document).ready(function() {
  settings.init();
});

var settings = function() {
  s = {};
  s.serializeFns = [];

  s.addSerializeFn = function(fn) {
    s.serializeFns.push(fn);
  };

  s.getSerializeFns = function() {
    return s.serializeFns;
  }

  s.init = function() {
    initForm();
  }

  var initForm = function() {
    $('#js-settings-form').submit(function(e) {
      e.preventDefault();
      var $form = $(this);
      var url = $form.attr('action'),
        formData = [],
        serializeFns = settings.getSerializeFns();
      $.each(serializeFns, function() {
        formData = formData.concat(this());
      });
      $.ajax({
        url: url,
        method: 'POST',
        data: JSON.stringify(formData),
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
        window.location.replace('/settings')
      }).fail(function(jqXHR, textStatus, errorThrown) {
        var $error = $('.js-error').first().empty();
        $('#js-errors').html($error);
        $error.append(jqXHR.responseText).removeClass('hidden');
      });
    });
  };

  return s;
}();
