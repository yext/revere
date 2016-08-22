$(document).ready(function() {
  emailTarget.init();
});

var emailTarget = function() {
  var e = {};

  e.init = function() {
    initRemove();
    removeEmail();
    addEmail();
    addSerializeFn();
  };

  var initRemove = function() {
    $.each($('.js-trigger'), function(i) {
      var $emails = $(this).find('.js-email-address');
      $.each($(this).find('.js-email-address'), function(j) {
        if (j === 0) {
          return;
        }
        $(this).find('button')
          .removeClass('js-add-email')
          .addClass('js-remove-email')
          .text('-');
      });
    });
  };

  var removeEmail = function() {
    $(document.body).on('click', '.js-remove-email', function(e) {
      e.preventDefault();
      $(this).parents('.js-email-address').remove();
    });
  };

  var addEmail = function() {
    $(document.body).on('click', '.js-add-email', function(e) {
      e.preventDefault();
      var $emailField = $('.js-email-address').first().clone();
      $emailField.find('input[type="text"]').val('');
      $emailField.find('button')
        .removeClass('js-add-email')
        .addClass('js-remove-email')
        .text('-');
      $emailField.appendTo($(this).parents('.js-target'));
    });
  };

  var addSerializeFn = function() {
    targets.addSerializeFn($('#js-email-target-type').val(), function(target) {
      var emailInputs = target.find(':input').serializeObject();

      // Deal with having multiple emails
      var emails = [];
      var emailTo = emailInputs.To;
      var replyTo = emailInputs.ReplyTo;
      if (typeof emailTo === 'string') {
        emails.push({'To': emailTo, 'ReplyTo': replyTo});
      } else {
        for (var i = 0; i < emailTo.length; i++) {
          emails.push({'To': emailTo[i], 'ReplyTo': replyTo[i]});
        }
      }
      return JSON.stringify({'Addresses':emails});
    });
  };

  return e;
}();
