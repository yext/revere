$(document).ready(function() {
  emailTarget.init();
});

var emailTarget = function() {
  var e = {};

  e.init = function() {
    addSerializeFn();
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

  return e;
}();
