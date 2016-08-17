$(document).ready(function() {
  emailTarget.init();
});

var emailTarget = function() {
  var e = {};

  e.init = function() {
    addSerializeFn();
    initTags();
  };

  var addSerializeFn = function() {
    targets.addSerializeFn($('#js-email-target-type').val(), function(target) {
      var emailInputs = target.find(':input').serializeObject();

      var emails = [];
      var emailTo = emailInputs.To.split(',');
      var replyTo = emailInputs.ReplyTo.split(',');
      for (var i = 0; i < Math.max(emailTo.length, replyTo.length); i++) {
        var ithEmailTo = i < emailTo.length ? emailTo[i].replace(/^\s+|\s+$/g,'') : '';
        var ithReplyTo = i < replyTo.length ? replyTo[i].replace(/^\s+|\s+$/g,'') : '';
        emails.push({'To': ithEmailTo, 'ReplyTo': ithReplyTo});
      }
      return JSON.stringify({'Addresses':emails});
    });
  };

  var initTags = function() {
    var tagConfig = {
      tagClass: 'label label-primary',
      trimValue: true
    };

    $('.js-trigger').find('.js-email-field:visible').tagsinput(tagConfig);

    // Initialize tags for new triggers
    MutationObserver = window.MutationObserver || window.WebKitMutationObserver;
    var target = document.querySelector('#triggers'),
      config = {childList: true},
      observer = new MutationObserver(function(mutations) {
      mutations.forEach(function(mutation) {
        [].slice.call(mutation.addedNodes).forEach(function(trigger) {
          $(trigger).find('.js-email-field').tagsinput(tagConfig);
        });
      });
    });
    observer.observe(target, config);
  };

  return e;
}();
