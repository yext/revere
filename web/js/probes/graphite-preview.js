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
  }

  var initEventHandlers = function() {
    previewGraphite();
    disableUnusedPreviewPeriod();
  }

  var previewGraphite = function() {
    $('#js-preview-btn').click(function(e) {
      e.preventDefault();
      var gtFields = $('#js-graphite-threshold').find(':input').serializeObject();
      var previewFields = $('#js-preview').find(':input').serializeObject();
      url = getGraphitePreviewUrl(
        getGraphiteBaseUrl(gtFields),
        getGraphiteTargets(gtFields),
        getGraphitePreviewPeriod(previewFields)
      );
      $('#js-preview-img').attr('src', url);
    });
  };

  var disableUnusedPreviewPeriod = function() {
    $(document.body).on('change','input[type=radio][name=previewPeriod]', function() {
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
    return gtFields['url'];
  }

  var getGraphiteTargets = function(gtFields) {
    return [
      gtFields['expression'],
      getThresholdTargetExpression(gtFields['warningThreshold'], 'warning', 'green'),
      getThresholdTargetExpression(gtFields['errorThreshold'], 'error', 'orange'),
      getThresholdTargetExpression(gtFields['criticalThreshold'], 'critical', 'red')
    ];
  }

  var getThresholdTargetExpression = function(threshold, label, color) {
    return 'threshold(' + threshold + ',"' + label + '","' + color + '")';
  }

  var getGraphitePreviewPeriod = function(previewFields) {
    var previewPeriod = {};
    if (previewFields['previewPeriod'] === 'last') {
      previewPeriod['from'] = '-' + previewFields['lastPeriod'] + previewFields['lastPeriodType'];
    } else if (previewFields['previewPeriod'] === 'range') {
      var fromDate = $('#js-datetimepicker-from').data().date,
        untilDate = $('#js-datetimepicker-until').data().date;

      previewPeriod['from'] = moment(fromDate).format('HH:mm_YYYYMMDD');
      previewPeriod['until'] = moment(untilDate).format('HH:mm_YYYYMMDD');
    }
    return previewPeriod;
  }

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
  }

  return gp;
}();
