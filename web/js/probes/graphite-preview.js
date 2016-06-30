$(document).ready(function() {
  graphitePreview.init();
});

var graphitePreview = function() {
  var gp = {};

  gp.init = function() {
    initDtps();
    initEventHandlers();
  };

  var initDtps = function() {
    var $fromDtp = $('#js-datetimepicker-from'),
      $untilDtp = $('#js-datetimepicker-until');

    $fromDtp.datetimepicker({
      format: revere.displayDateTimeFormat(),
      defaultDate: moment().subtract(1, 'days')
    });
    $untilDtp.datetimepicker({
      format: revere.displayDateTimeFormat(),
      useCurrent: false,
      defaultDate: moment()
    });
    $fromDtp.on('dp.change', function(e) {
      $untilDtp.data('DateTimePicker').minDate(e.date);
    });
    $untilDtp.on('dp.change', function(e) {
      $fromDtp.data('DateTimePicker').maxDate(e.date);
    });
  };

  var initEventHandlers = function() {
    previewGraphite();
    disableUnusedPreviewPeriod();
  };

  var previewGraphite = function() {
    var $previewImg = $('#js-preview-img'),
      $previewBtn = $('#js-preview-btn'),
      $previewError = $('#js-preview-error');

    $previewImg.on('load', function() {
      $previewError.addClass('hidden');
      $previewImg.removeClass('hidden');
      $previewBtn.button('reset');
    });

    $previewImg.on('error', function() {
      $previewImg.addClass('hidden');
      $previewError.removeClass('hidden');
      $previewBtn.button('reset');
    });

    $previewBtn.click(function(e) {
      e.preventDefault();
      $(this).button('loading');
      var gtFields = $('#js-graphite-threshold').find(':input').serializeObject();
      var previewFields = $('#js-preview').find(':input').serializeObject();
      url = getGraphitePreviewUrl(
        getGraphiteBaseUrl(gtFields),
        getGraphiteTargets(gtFields),
        getGraphitePreviewPeriod(previewFields)
      );
      $previewImg.attr('src', url);
    });
  };

  var disableUnusedPreviewPeriod = function() {
    $(document.body).on('change','input[type=radio][name=PreviewPeriod]', function() {
      if ($(this).val() === 'last') {
        $('.js-range-period').prop('disabled', true);
        $('.js-last-period').prop('disabled', false);
      } else if ($(this).val() === 'range') {
        $('.js-last-period').prop('disabled', true);
        $('.js-range-period').prop('disabled', false);
      }
    });
  };

  var getGraphiteBaseUrl = function(gtFields) {
    return gtFields['URL'];
  };

  var getGraphiteTargets = function(gtFields) {
    return [
      getDataTargetExpression(gtFields['Expression'], gtFields['TriggerIf']),
      getThresholdTargetExpression(gtFields['Warning'], 'warning', 'orange'),
      getThresholdTargetExpression(gtFields['Error'], 'error', 'red'),
      getThresholdTargetExpression(gtFields['Critical'], 'critical', 'black')
    ];
  };

  var getDataTargetExpression = function(expression, triggerIf) {
    var numSeries = 3;
    if (triggerIf == '>' || triggerIf == '>=') {
      return 'highestMax(' + expression + ',' + numSeries + ')';
    } else if (triggerIf == '<' || triggerIf == '<=') {
      // XXX: lowestMin does not exist in graphite
      // Scale by -1 use highest max and then scale back. Also alias out scale funcs
      return 'aliasSub(scale(highestMax(scale(' +
        expression + ',-1),' + numSeries + '),-1),"^(?:scale\\(){2}(.+)(?:,-1\\)){2}$","\\1")';
    } else {
      return expression;
    }
  };

  var getThresholdTargetExpression = function(threshold, label, color) {
    return 'threshold(' + threshold + ',"' + label + '","' + color + '")';
  };

  var getGraphitePreviewPeriod = function(previewFields) {
    var previewPeriod = {};
    if (previewFields['PreviewPeriod'] === 'last') {
      previewPeriod['from'] = '-' + previewFields['LastPeriod'] + previewFields['LastPeriodType'];
    } else if (previewFields['PreviewPeriod'] === 'range') {
      var fromDate = $('#js-datetimepicker-from').data().date,
        untilDate = $('#js-datetimepicker-until').data().date;

      previewPeriod['from'] = moment(fromDate).format('HH:mm_YYYYMMDD');
      previewPeriod['until'] = moment(untilDate).format('HH:mm_YYYYMMDD');
    }
    return previewPeriod;
  };

  var getGraphitePreviewUrl = function(baseUrl, targets, previewPeriod) {
    params = {
      // Cache buster
      '_salt': Date.now(),
      'height': 200,
      'width': 600,
      'target': targets
    }
    params = $.extend(params, previewPeriod);
    return '//' + baseUrl + '/render/?' + $.param(params);
  };

  return gp;
}();
