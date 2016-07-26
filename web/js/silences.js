$(document).ready(function() {
  silencesEdit.init();
});

var silencesEdit = function() {
  var s = {};
  var defaultDtpSettings = {
    useCurrent: false,
    sideBySide: true,
    format: revere.displayDateTimeFormat()
  };
  var defaultStartEndOffset = 60 * 60; // 1 hour offset in unix timestamp

  var $startDtp = $('.js-datetimepicker-start'),
    $endDtp = $('.js-datetimepicker-end');

  // Determines whether the request is creating or editing a silence
  var isNew = $startDtp.data('time') === revere.goTimeZero();

  s.init = function() {
    initDtps();
    initForm();
    initEndNow();
    initSilenceBounds();
  };

  var initSilenceBounds = function() {
    var startDtp = $startDtp.data('DateTimePicker'),
      endDtp = $endDtp.data('DateTimePicker');

    if (isNew) {
      $('#js-start-now, #js-end-duration').prop('checked', true);
    } else {
      $('#js-start-dtp, #js-end-dtp').prop('checked', true);
    }
  };

  var initDtps = function() {
    var nowEpoch = moment().seconds(0).format(revere.modelDateTimeFormat()),
      modelStart = $startDtp.data('time'),
      modelEnd = $endDtp.data('time');

    var startEpoch = defaultToNow(nowEpoch, modelStart),
      endEpoch = defaultToNow(addDefaultOffsetToEpoch(nowEpoch), modelEnd);

    var now = epochToLocalTimeString(nowEpoch),
      start = epochToLocalTimeString(startEpoch),
      end = epochToLocalTimeString(endEpoch);

    setStartDtp(now, start);
    setEndDtp(start, end);

    var startDtp = $startDtp.data('DateTimePicker'),
      endDtp = $endDtp.data('DateTimePicker');

    $startDtp.on('dp.change', function(e) {
      var now = moment(),
        start = e.date,
        end = endDtp.date();

      if (end.isSameOrBefore(start)) {
        end = moment(start);
        end.add(1, 'hour');
      }

      setStartDtp(now, start);
      setEndDtp(start, end.format(revere.displayDateTimeFormat()));
    });

    disableInvalidBoundFields(startDtp, endDtp);
  };

  var initForm = function() {
    $('.js-submit-btn').click(function(e) {
      e.preventDefault();
      saveSilence();
    });
  };

  var initEndNow = function() {
    var now = moment(),
      endDtp = $endDtp.data('DateTimePicker');

    var $startNow = $('#js-start-now');
    $('#js-end-silence').click(function(e) {
      e.preventDefault();
      if ($startNow.is(':enabled')) {
        $startNow.prop('checked', true);
      }
      $('#js-end-dtp').prop('checked', true);

      endDtp.minDate(now);
      endDtp.date(now);
      saveSilence();
    });
  };

  var saveSilence = function() {
    var $invalidInput = $('.js-invalid-input')
      .addClass('hidden').empty();
    var $validInput = $('.js-valid-input')
      .addClass('hidden').empty();
    var $serverError = $('.js-server-error')
      .addClass('hidden').empty();

    var startDtp = $startDtp.data('DateTimePicker'),
      endDtp = $endDtp.data('DateTimePicker');

    var startMoment = ($('.js-start-type:checked').val() === 'now') ?
      moment() : startDtp.date();

    var endMoment = ($('.js-end-type:checked').val() == 'duration') ?
      getEndMomentFromDuration(startMoment) : endDtp.date();

    // Prevents serialization from picking up fields
    disableAllBoundFields();

    var data = getSilenceData(),
      id = data['SilenceId'];
    data.Start = localTimeToUtc(startMoment);
    data.End = localTimeToUtc(endMoment);
    $.ajax({
      method: 'POST',
      url: '/silences/'+ id + '/edit',
      data: JSON.stringify(data),
      contentType: 'application/json; charset=UTF-8'
    }).success(function(d) {
      if (d.errors) {
        return revere.showErrors(d.errors);
      } else {
        var timeout = window.setTimeout(function() {
          window.location.replace('/silences/'+d.id);
        }, 200);
      }
    }).fail(function(jqXHR, status) {
      $serverError.append('<h4>Server Error</h4>');
      $serverError.append('<p>' + jqXHR.responseText + '</p>');
      $serverError.removeClass('hidden');
      enableBoundFields();
      disableInvalidBoundFields(startDtp, endDtp);
    });
  };

  var setStartDtp = function(now, start) {
    $startDtp.datetimepicker(defaultDtpSettings);
    var $dtpObj = $startDtp.data('DateTimePicker');
    $dtpObj.defaultDate(start);
    $dtpObj.minDate(now)
    $dtpObj.date(start);
  };

  var setEndDtp = function(start, end) {
    $endDtp.datetimepicker(defaultDtpSettings);
    var $dtpObj = $endDtp.data('DateTimePicker');
    $dtpObj.defaultDate(end);
    $dtpObj.minDate(start);
    $dtpObj.maxDate(moment(start).add(2, 'week').format(revere.displayDateTimeFormat()));
    $dtpObj.date(end);
  };

  var getEndMomentFromDuration = function(startMoment) {
    var $duration = $('input[name="duration"]'),
      $durationType = $('select[name="durationType"]');
    $duration.prop('disabled', true);
    $durationType.prop('disabled', true);
    return moment(startMoment).add($duration.val(), $durationType.val());
  };

  var disableAllBoundFields = function() {
    $('#silence-bounds').find(':input').prop('disabled', true);
  };

  var disableInvalidBoundFields = function(startDtp, endDtp) {
    if (isNew) {
      return;
    }
    var now = moment().seconds(0);
    if (moment(startDtp.date()).isBefore(now)) {
      if (moment(endDtp.date()).isSameOrBefore(now)) {
        disableAllBoundFields();
        return;
      }
      $('.js-start-type').attr('disabled', true);
      startDtp.disable();
    }
  };

  var enableBoundFields = function() {
    $('#silence-bounds').find(':input').prop('disabled', false);
  };

  var defaultToNow = function(now, time) {
    return time === revere.goTimeZero() ? now : time;
  };

  var epochToLocalTimeString = function(epoch) {
    return moment.unix(epoch).tz(revere.localTimeZone()).format(revere.displayDateTimeFormat());
  };

  var addDefaultOffsetToEpoch = function(epoch) {
    return (parseInt(epoch) + parseInt(defaultStartEndOffset)).toString();
  };

  var getSilenceData = function() {
    return $('#js-silence-info').find(':input').serializeObject();
  };

  var localTimeToUtc = function(time) {
    return moment.tz(time.format(revere.displayDateTimeFormat()), revere.localTimeZone()).tz(revere.serverTimeZone());
  };

  return s;
}();

