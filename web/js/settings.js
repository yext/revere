var outgoingEmails = function() {
  oe = {};

  oe.init = function() {
      initForm();
  }

  var initForm = function() {
    $('.submit-settings-button').click(function() {
      var mySection = $(this).parents('.js-setting-section');
      var postJson = mySection.find('.js-setting-input').serializeObject();
      $.ajax({
        url: mySection.data('endpoint'),
        method: 'POST',
        contentType: 'application/json; charset=UTF-8',
        context: this,
        data: JSON.stringify(postJson)
      }).success(function(response) {
        if (response.errors) {
          var $error = $('.error').first().empty();
          $('#errors').html($error);
          $.each(response.errors, function() {
            $error.append('<p>' + this + '<p/>').removeClass('hidden');
          });
        } else {
          location.reload();
        }
      }).fail(function(jqXHR, textStatus, errorThrown) {
        var $error = $('.error').first().empty();
        $('#errors').html($error);
        $error.append(jqXHR.responseText).removeClass('hidden');
      });
    });
  }
  
  return oe;
}();
