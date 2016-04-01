$(document).ready(function() {
  subprobesIndex.init();
});

var subprobesIndex = function() {
  var s = {};

  s.init = function() {
    initEnteredStates();
  };

  var initEnteredStates = function () {
    $('.js-subprobe-entered-state').each(function(i) {
      var $this = $(this);
      var title = moment($this.prop('title'), 'YYYY-MM-DD- HH:mm:ss ZZ').format('YYYY-MM-DD HH:mm:ss UTCZZ');
      $this.attr('title', title);
    });
    $('.js-subprobe-entered-state').tooltip({container: 'body'});
  };

  return s;
}();