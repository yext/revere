$('#js-preview-btn').on('btn-loaded', function() {
  $('#js-preview-btn').click()
});

$(document).ready(function(){
  deleteSubprobe();
});

var deleteSubprobe = function() {
  $('#delete').click(function() {
    $.ajax({
      url: window.location.pathname + '/delete',
      method: 'DELETE',
      contentType: 'application/json; charset=UTF-8'
    }).success(function (response) {
      if (response.errors) {
        return revere.showErrors(response.errors);
      }
      window.location.href = './';
    });
  });
};
