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
