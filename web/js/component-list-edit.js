var componentList = function() {
  return function(componentName) {
    var cl = {};

    cl.init = function() {
      initComponents();
    };

    cl.getData = function() {
      var data = [];
      $.each($.merge($('.js-' + componentName),
          $('.js-new-' + componentName + ':visible')), function() {
        component = $(this).find(':input').serializeObject();
        data.push(component);
      });
      return data;
    };

    var initComponents = function() {
      var showOption = function(id, $selectBoxes) {
        $.each($selectBoxes, function() {
          $(this).find('option[value="' + id + '"]').removeClass('hidden');
        });
      };

      var hideOption = function(id, $selectBoxes) {
        if (id === '') {
          return;
        }
        $.each($selectBoxes, function() {
          $(this).find('option[value="' + id + '"]').addClass('hidden');
        });
      };

      var $baseAddComponent = $('.js-new-' + componentName).first();

      $.each($('.js-' + componentName), function() {
        var id = $(this).find('input.js-id').val();
        hideOption(id, $('.js-' + componentName + '-name'));
      });

      $('#js-add-' + componentName).click(function(e) {
        e.preventDefault();
        var $existingComponents = $('.js-' + componentName),
          $newComponents = $('.js-new-' + componentName),
          $newComponent = $baseAddComponent.clone();

        $('.js-empty-' + componentName).addClass('hidden');
        if ($newComponents.length > 1 || $existingComponents.length === 0) {
          $newComponent.insertAfter($newComponents.last());
        } else {
          $newComponent.insertAfter($existingComponents.last());
        }
        $newComponent.removeClass('hidden');
      });

      $(document.body).on('click', '.js-remove-' + componentName, function(e) {
        e.preventDefault();
        $component = $(this).parents('.js-' + componentName);
        $component.hide();
        $component.find('input[name="Delete"]').prop('checked', true);
        showOption($component.find('input.js-id').val(), $('.js-' + componentName + '-name'));
      });

      $(document.body).on('click', '.js-remove-new-' + componentName, function(e) {
        e.preventDefault();
        $this = $(this),
          $newComponent = $this.parents('.js-new-' + componentName),
          $selectedOption = $newComponent.find('.js-' + componentName + '-name option:selected');

        showOption($selectedOption.val(), $('.js-' + componentName + '-name'));
        $newComponent.remove();
      });

      $(document.body).on('change', '.js-' + componentName + '-name', function() {
        var $this = $(this),
          $selectedOption = $this.find('option:selected'),
          $otherSelectBoxes = $('.js-' + componentName + '-name').not($this);

        var previousId = $this.data('id'),
          description = $selectedOption.data('description');

        $this.parents('.js-new-' + componentName).find('.description').text(description);
        $this.data('id', $selectedOption.val());

        showOption(previousId, $otherSelectBoxes);
        hideOption($selectedOption.val(), $otherSelectBoxes);
      });
    };

    return cl;
  };
}();
