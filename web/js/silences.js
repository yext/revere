var silenceDtp = function() {
  var sdtp = {};
  var defaultDtpSettings = {
    useCurrent: false,
    sideBySide: true
  };

  sdtp.init = function(format, id) {
    defaultDtpSettings['format'] = format
    initDtps();
    initForm(id);
  };

  var initDtps = function() {
    var $startDtp = $('.js-datetimepicker-start'),
      $endDtp = $('.js-datetimepicker-end');
    var now = moment(),
      startTime = $startDtp.data('time'),
      endTime = $endDtp.data('time'),
      start = getStartMoment(now, startTime),
      end = getEndMoment(now, endTime);

    setStartDtp($startDtp, now, start);
    setEndDtp($endDtp, now, start, end);

    $startDtp.on('dp.change', function(e) {
      var now = moment(),
        start = moment(e.date),
        end = moment($endDtp.data('DateTimePicker').date());

      if (end.isBefore(start)) {
        end = moment(start);
        end.add(1, 'hour');
      }

      setStartDtp($startDtp, now, start);
      setEndDtp($endDtp, now, start, end);
    });

    if (start.isBefore(now)) {
      $startDtp.data('DateTimePicker').disable();
      if (end.isBefore(now)) {
        $endDtp.data('DateTimePicker').disable();
        return;
      }
    }

  };

  var initForm = function(id) {
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
      var data = getSilenceData();
      data.start = $startDtp.date().format();
      data.end = $endDtp.date().format();
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
    $dtpObj.maxDate(moment(start).add(2, 'week'));
    $dtpObj.date(end);
  }

  var getStartMoment = function(now, start) {
    return start === '' ? moment(now) : moment(start, defaultDtpSettings.format);
  };

  var getEndMoment = function(now, end) {
    return end === '' ? moment(now).add(1, 'hour') : moment(end, defaultDtpSettings.format);
  };

  var getSilenceData = function() {
    return $('#js-silence-form').find(':input').serializeObject();
  }

  return sdtp;
}();

