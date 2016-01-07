$(document).ready(function() {
  // TODO(psingh): Put in appropriately named fn and remove comment
  // Have a remove option when there are multiple emails in a trigger
  $.each($('.trigger'), function(i) {
    var $emails = $(this).find('.email-address');
    $.each($(this).find('.email-address'), function(j) {
      if (j == 0) {
        return;
      }
      $(this).find('button')
        .removeClass('add-email')
        .addClass('remove-email')
        .text('-');
    });
  });

  $(document.body).on('click', '.add-email', function(e) {
    e.preventDefault();
    var $emailField = $('.email-address').first().clone();
    $emailField.find('input[type="text"]').val('');
    $emailField.find('button')
      .removeClass('add-email')
      .addClass('remove-email')
      .text('-');
    $emailField.insertAfter($(this).parents('.email-address'));
  });

  $(document.body).on('click', '.remove-email', function(e) {
    e.preventDefault();
    $(this).parents('.email-address').remove();
  });

  targets.addSerializeFn($('#email-target-type').val(), function(target) {
    var emailInputs = target.find(':input').serializeObject();

    // Deal with having multiple emails
    var emails = [];
    var emailTo = emailInputs.emailTo;
    var replyTo = emailInputs.replyTo;
    if (typeof emailTo === 'string') {
      emails.push({'emailTo': emailTo, 'replyTo': replyTo});
    } else {
      for (var i = 0; i < emailTo.length; i++) {
        emails.push({'emailTo': emailTo[i], 'replyTo': replyTo[i]});
      }
    }
    return JSON.stringify({'emails':emails});
  });
});

