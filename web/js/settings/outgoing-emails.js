$(document).ready(function() {
  settings.addSerializeFn(outgoingEmails.getData);
});


var outgoingEmails = function() {
  var oe = {};

  oe.getData = function() {
    var data = [];
    $.each($('.js-outgoing-email'), function() {
      var serialized = $(this).find(':input.required').serializeObject();
      var json = $(this).find(':input.json').serializeObject();
      $.extend(serialized, {'setting': JSON.stringify(json)});
      data.push(serialized);
    });
    return data;
  };

  return oe;
}();
