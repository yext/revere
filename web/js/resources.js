$(document).ready(function() {
  resources.init();
});

var resources = function() {
  var dsi = {};
  dsi.sourceFunctions = [];

  dsi.addSourceFunction = function(fn) {
    dsi.sourceFunctions.push(fn);
  };

  dsi.getSourceFunctions = function() {
    return dsi.sourceFunctions;
  };

  var initForm = function() {
    $('#js-resources-form').submit(function(e) {
      e.preventDefault();
      var $form = $(this);
      var url = $form.attr('action'),
        formData = [],
        sourceFns = resources.getSourceFunctions();
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
        window.location.replace('/resources');
      }).fail(function(jqXHR, textStatus, errorThrown) {
        revere.showErrors([jqXHR.responseText || textStatus]);
      });
    });
  };

  var clearInputs = function(newField) {
    newField.find('input[type="text"]').val('');
    newField.find('input[name="ResourceID"]').val(0);
    newField.find('input[name="Delete"]').val(false);
  };

  var initDeleteButtons = function() {
    $(document.body).on('click', '.js-remove-resource', function(e) {
      e.preventDefault();
      $resource = $(this).parents('.js-resource');
      var id = $resource.find('input[name="ResourceID"]').val();
      if(id == '0'){
        $resource.remove();
      } else if (resourcesLeft() === 1) {
        clearInputs($resource);
      } else {
        $resource.find('input[name="Delete"]').prop('checked', true);
        $resource.addClass('hidden');
      }
    });
  }

  var resourcesLeft = function() {
    var type = $(this).data('sourceref');
    return $('.' + type).not('hidden').length;
  }

  var sortAndArrangeResources = function() {
      // Get resources as they are displayed
      var $resources = $('.js-resource');
      idToHtmlMap = {};
      var id;
      $.each($resources, function(_index, sourceHtml) {
          id = parseInt($(sourceHtml).find('input[name="ResourceType"]').val());
          if (idToHtmlMap[id]) {
              idToHtmlMap[id].push(sourceHtml);
          } else {
              idToHtmlMap[id] = [];
              idToHtmlMap[id].push(sourceHtml);
          }
      });
      var $typeDivs = $('.js-resource-type');

      $('#js-resource-list').remove();
      $.each($typeDivs, function(_index, div) {
          idStr = $(div).attr('js-resource-type');
          id = parseInt(idStr);
          $button = $(div).find('button');
          if (idToHtmlMap[id]) {
            for (var i = 0; i < idToHtmlMap[id].length; i++) {
              $button.before(idToHtmlMap[id][i]);
            }
          }
          $button.click(function() {
              $.ajax('resourcetype/'+idStr, {
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
    sortAndArrangeResources();
  };
  return dsi;
}();
