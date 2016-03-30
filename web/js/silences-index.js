$(document).ready(function() {
  silencesIndex.init();
});

var silencesIndex = function () {
  var s = {};

  s.init = function() {
    var $silences = $('.js-silence'),
      $futureSilenceDiv = $('.js-future-silences'),
      $currentSilenceDiv = $('.js-current-silences'),
      $pastSilenceDiv = $('.js-past-silences');

    $silences.sort(compareSilences);

    var now = moment().seconds(0);

    $silences.each(function() {
      var $silence = $(this);

      var $start = getStartElement($silence),
        $end = getEndElement($silence);

      var start = getStartMoment($start),
        end = getEndMoment($end);

      setText($start, start);
      setText($end, end);

      if (now.isBefore(start)) {
        $silence.appendTo($futureSilenceDiv);
      } else if (now.isBefore(end)) {
        $silence.appendTo($currentSilenceDiv);
      } else {
        $silence.appendTo($pastSilenceDiv);
      }
    });

    $futureSilences = $futureSilenceDiv.children('.js-silence')
    $currentSilences = $currentSilenceDiv.children('.js-silence')
    $pastSilences = $pastSilenceDiv.children('.js-silence')

    if ($futureSilences.length > 0) {
      $futureSilenceDiv.removeClass('hidden');
      $futureSilences.removeClass('hidden');
    }

    if ($currentSilences.length > 0) {
      $currentSilenceDiv.removeClass('hidden');
      $currentSilences.removeClass('hidden');
    }

    if ($pastSilences.length > 0) {
      $pastSilenceDiv.removeClass('hidden');
      $pastSilences.removeClass('hidden');
    }
  };

  var getStartElement = function ($silence) {
    return $silence.find('.js-silence-start');
  }

  var getEndElement = function ($silence) {
    return $silence.find('.js-silence-end');
  }

  var getStartMoment = function ($silenceStart) {
    return moment.unix($silenceStart.data('time'));
  };

  var getEndMoment = function ($silenceEnd) {
    return moment.unix($silenceEnd.data('time'));
  };

  var setText = function ($silenceTimeDiv, silenceMoment) {
    $silenceTimeDiv.text(silenceMoment.format(revere.displayDateTimeFormat()))
  }

  var compareSilences = function (silence1, silence2) {
    var start1 = getStartMoment(getStartElement($(silence1))),
      start2 = getStartMoment(getStartElement($(silence2)));

    if (start1.isBefore(start2)){
      return 1;
    } else if (start2.isBefore(start1)) {
      return -1;
    } else {
      return 0;
    }
  }

  return s
}();
