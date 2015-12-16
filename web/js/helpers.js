$.fn.serializeObject = function() {
  var o = {};
  var a = this.serializeArray();

  // Only keep form inputs that serializeArray would
  var that = this;
  $.each(that, function(i) {
    if (!$(this).is('[name]') || $(this).prop('disabled') ||
        ($(this).is(':checkbox') && !$(this).is(':checked'))) {
      that.splice(i, 1);
    }
  });
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
