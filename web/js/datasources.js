$(document).ready(function() {
  datasources.init();
});

var datasources = function() {
  var dsi = {};
  dsi.sourceFunctions = [];

  dsi.addSourceFunction = function(fn) {
    dsi.sourceFunctions.push(fn);
  };

  dsi.getSourceFunctions = function() {
    return dsi.sourceFunctions;
  };

  var initForm = function() {
    $('#js-datasources-form').submit(function(e) {
      e.preventDefault();
      var $form = $(this);
      var url = $form.attr('action'),
        formData = [],
        sourceFns = datasources.getSourceFunctions();
      $.each(sourceFns, function() {
        formData = formData.concat(this());
      });
      $.ajax({
        url: url,
        method: 'POST',
        data: JSON.stringify(formData),
        contentType: 'application/json; charset=UTF-8'
      }).success(function(response) {
        if (response.errors) {
          return revere.showErrors(response.errors);
        }
        window.location.replace('/datasources');
      }).fail(function(jqXHR, textStatus, errorThrown) {
        revere.showErrors([jqXHR.responseText || textStatus]);
      });
    });
  };

  var clearInputs = function(newField) {
    newField.find('input[type="text"]').val('');
    newField.find('input[name="SourceID"]').val(0);
    newField.find('input[name="Delete"]').val(false);
  };

  var initDeleteButtons = function() {
    $(document.body).on('click', '.js-remove-datasource', function(e) {
      e.preventDefault();
      $dataSource = $(this).parents('.js-datasource');
      var id = $dataSource.find('input[name="SourceID"]').val();
      if(id == '0'){
        $dataSource.remove();
      } else if (datasourcesLeft === 1) {
        clearInputs($dataSource);
      } else {
        $dataSource.find('input[name="Delete"]').prop('checked', true);
        $dataSource.addClass('hidden');
      }
    });
  }

  var datasourcesLeft = function() {
    var type = $(this).data('sourceref');
    return $('.' + type).not('hidden').length;
  }

  var sortAndArrangeDataSources = function() {
      // Get datasources as they are displayed
      var $datasources = $('.js-datasource');
      idToHtmlMap = {};
      var id;
      $.each($datasources, function(_index, sourceHtml) {
          id = parseInt($(sourceHtml).find('input[name="SourceType"]').val());
          if (idToHtmlMap[id]) {
              idToHtmlMap[id].push(sourceHtml);
          } else {
              idToHtmlMap[id] = [];
              idToHtmlMap[id].push(sourceHtml);
          }
      });
      var $typeDivs = $('.js-source-type');

      $('#js-sources').remove();
      $.each($typeDivs, function(_index, div) {
          idStr = $(div).attr('js-source-type');
          id = parseInt(idStr);
          $button = $(div).find('button');
          if (idToHtmlMap[id]) {
            for (var i = 0; i < idToHtmlMap[id].length; i++) {
              $button.before(idToHtmlMap[id][i]);
            }
          }
          $button.click(function() {
              $.ajax('datasourcetype/'+idStr, {
                  success: function(data) {
                    $button.before(data['template']);
                  }
              });
          });
      });
  };

  dsi.init = function() {
    initForm();
    initDeleteButtons();
    sortAndArrangeDataSources();
  };
  return dsi;
}();
