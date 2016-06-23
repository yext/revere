var targets = function() {
  var t = {};
  var fns = {};

  t.addSerializeFn = function(targetType, fn) {
    fns[targetType] = fn;
  };

  t.getSerializeFn = function(targetType) {
    return fns[targetType];
  };

  return t;
}();

var triggerEdit = function() {
  var te = {};

  te.getData = function(jstrigger) {
    var triggerOptions = $(jstrigger).find('.js-trigger-options :input').serializeObject(),
      targetFn = targets.getSerializeFn(triggerOptions['TargetType']),
      target;
    if (targetFn !== undefined) {
      target = targetFn($(jstrigger).find('.js-target'));
    }

    if (target === undefined) {
      target = JSON.stringify($(jstrigger).find('.js-target :input').serializeObject());
    }
    return $.extend(triggerOptions, {'TargetParams':target});
  };

  return te;
}();
