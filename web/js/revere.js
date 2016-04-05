var revere = function() {
  r = {};

  r.displayDateTimeFormat = function() {
    return 'YYYY-MM-DD HH:mm';
  };

  r.modelDateTimeFormat = function () {
    return 'X';
  };

  r.goTimeZero = function () {
    return -62135596800;
  };

  r.serverTimeZone = function() {
    return 'UTC';
  };

  r.localTimeZone = function() {
    return moment.tz.guess();
  };

  r.getParameterByName = function(name) {
    var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
    return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
  };

  return r;
}();

$.fn.serializeObject = function() {
  var o = {};
  var a = this.serializeArray();

  // Only keep form inputs that serializeArray would
  var that = this;
  for(var i = 0; i < that.length; i++) {
    if (!$(that[i]).is('[name]') || $(that[i]).is(':disabled') ||
        ($(that[i]).is(':checkbox') && !$(that[i]).is(':checked')) ||
        ($(that[i]).is(':radio') && !$(that[i]).is(':checked'))) {
      that.splice(i, 1);
      i--;
    }
  }
  $.each(a, function(i) {
    var value = this.value || '';
    if (that[i].dataset && that[i].dataset.jsonType) {
        // Deal with checkboxes, set defaults, parse everything else
        if (value === 'on') {
          value = true;
        } else if (this.value === '' && that[i].dataset.jsonType == 'Number') {
          value = 0;
        } else if (this.value === '' && that[i].dataset.jsonType == 'Boolean') {
          value = false;
        } else {
          value = JSON.parse(value);
        }
    }
    if (o[this.name] !== undefined) {
      if (!o[this.name].push) {
        o[this.name] = [o[this.name]];
      }
      o[this.name].push(value);
    } else {
      o[this.name] = value;
    }
  });
  return o;
};
