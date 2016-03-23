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

  s.init = function() {
    initDtps();
    initForm();
  };

  var initDtps = function() {
    var $startDtp = $('.js-datetimepicker-start'),
      $endDtp = $('.js-datetimepicker-end');

    var nowEpoch = moment().seconds(0).format(revere.modelDateTimeFormat()),
      modelStart = $startDtp.data('time'),
      modelEnd = $endDtp.data('time');

    var startEpoch = defaultToNow(nowEpoch, modelStart),
      endEpoch = defaultToNow(addDefaultOffsetToEpoch(nowEpoch), modelEnd);

    var now = epochToLocalTimeString(nowEpoch),
      start = epochToLocalTimeString(startEpoch),
      end = epochToLocalTimeString(endEpoch);

    setStartDtp($startDtp, now, start);
    setEndDtp($endDtp, now, start, end);

    $startDtp.on('dp.change', function(e) {
      var now = moment(),
        start = e.date,
        end = $endDtp.data('DateTimePicker').date();

      if (end.isBefore(start)) {
        end = moment(start);
        end.add(1, 'hour');
      }

      setStartDtp($startDtp, now, start);
      setEndDtp($endDtp, now, start, end.format(revere.displayDateTimeFormat()));
    });

    if (moment(start).isBefore(now)) {
      $startDtp.data('DateTimePicker').disable();
      if (moment(end).isBefore(now)) {
        $endDtp.data('DateTimePicker').disable();
        return;
      }
    }

  };

  var initForm = function() {
    $('.js-submit-btn').click(function(e) {
      e.preventDefault();
      var $invalidInput = $('.js-invalid-input')
        .addClass('hidden').empty();
      var $validInput = $('.js-valid-input')
        .addClass('hidden').empty();
      var $serverError = $('.js-server-error')
        .addClass('hidden').empty();

      var $startDtp = $('.js-datetimepicker-start').data('DateTimePicker'),
        $endDtp = $('.js-datetimepicker-end').data('DateTimePicker');
      var startDtpMoment = $startDtp.date(),
        endDtpMoment = $endDtp.date();
      var data = getSilenceData();
      var id = data['id'];
      data.start = localTimeToUtc(startDtpMoment);
      data.end = localTimeToUtc(endDtpMoment);
      $.ajax({
        method: 'POST',
        url: '/silences/'+ id + '/edit',
        data: JSON.stringify(data),
        contentType: 'application/json; charset=UTF-8'
      }).success(function(d) {
        if (d.errors) {
          $.each(d.errors, function(i, v) {
            $invalidInput.append('<p>' + v + '</p>');
          });
          $invalidInput.removeClass('hidden');
        } else {
          var $validInput = $('.js-valid-input');
          if (id === 'new') {
            $validInput.append('<p>Successfully created silence</p>');
          } else {
            $validInput.append('<p>Successfully updated silence</p>');
          }
          $validInput.removeClass('hidden');
          var timeout = window.setTimeout(function() {
            window.location.replace('/silences/'+d.id);
          }, 200);
        }
      }).fail(function(jqXHR, status) {
        $serverError.append('<h4>Server Error</h4>');
        $serverError.append('<p>' + jqXHR.responseText + '</p>');
        $serverError.removeClass('hidden');
      });
    });
  };

  var setStartDtp = function($startDtp, now, start) {
    $startDtp.datetimepicker(defaultDtpSettings);
    var $dtpObj = $startDtp.data('DateTimePicker');
    $dtpObj.defaultDate(start);
    $dtpObj.minDate(now)
    $dtpObj.date(start);
  }

  var setEndDtp = function($endDtp, now, start, end) {
    $endDtp.datetimepicker(defaultDtpSettings);
    var $dtpObj = $endDtp.data('DateTimePicker');
    $dtpObj.defaultDate(end);
    $dtpObj.minDate(start);
    $dtpObj.maxDate(moment(start).add(2, 'week').format(revere.displayDateTimeFormat()));
    $dtpObj.date(end);
  }

  var defaultToNow = function(now, time) {
    return time === revere.goTimeZero() ? now : time;
  };

  var epochToLocalTimeString = function(epoch) {
    return moment.unix(epoch).tz(revere.localTimeZone()).format(revere.displayDateTimeFormat());
  };

  var addDefaultOffsetToEpoch = function(epoch) {
    return (parseInt(epoch) + parseInt(defaultStartEndOffset)).toString();
  }

  var getSilenceData = function() {
    return $('#js-silence-info').find(':input').serializeObject();
  }

  var localTimeToUtc = function(time) {
    return moment.tz(time.format(revere.displayDateTimeFormat()), revere.localTimeZone()).tz(revere.serverTimeZone());
  }

  return s;
}();

