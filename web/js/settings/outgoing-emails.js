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
      $.extend(serialized, {'SettingParams': JSON.stringify(json)});
      // TODO(fchen) SettingID and SettingType serialization should occur at the general setting level and shouldn't require somebody who is extending revere to implement this all over again
      data.push(serialized);
    });
    return data;
  };

  return oe;
}();
