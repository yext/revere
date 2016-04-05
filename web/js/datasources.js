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
                    var $error = $('.js-error').first().empty();
                    $('#js-errors').html($error);
                    $.each(response.errors, function() {
                        $error.append(this + '<br/>').removeClass('hidden');
                    });
                    return;
                }
                window.location.replace('/datasources')
            }).fail(function(jqXHR, textStatus, errorThrown) {
                var $error = $('.js-error').first().empty();
                $('#js-errors').html($error);
                $error.append(jqXHR.responseText).removeClass('hidden');
            });
        });
    };

    var clearInputs = function(newField) {
        newField.find('input[type="text"]').val('');
        newField.find('input[name="id"]').val(0);
    };

    var deleteFunc = function(e) { 
        var url = '/datasources/delete';
        $row = $(this).closest('.js-data-source');
        var id = $row.find('input[name="id"]').val();
        var data = {'id': parseInt(id)};
        $.ajax({
            url: url,
            method: 'POST',
            data: JSON.stringify(data),
            contentType: 'application/json; charset=UTF-8'
        }).success(function(response) {
            if ($row.parent().find('.js-data-source').length > 1) {
                $row.remove();
            } else {
                clearInputs($row);
            }
        }).fail(function(jqXHR, textStatus, errorThrown) {
            var $error = $('.js-error').first().empty();
            $('#js-errors').html($error);
            $error.append(jqXHR.responseText).removeClass('hidden');
        });
    };


    var initAddButtons = function() {
        $('.js-add-source').click(function(e) {
            e.preventDefault();
            var type = $(this).data('sourceref');
            newField = $('.' + type).first()
                .clone()
                .insertAfter('.' + type + ':last');
            clearInputs(newField)
            newField.find('.delete').click(deleteFunc);
        });
    };

    var initDeleteButtons = function() {
        $('.delete').click(deleteFunc)
    };

    dsi.init = function() {
        initForm();
        initAddButtons();
        initDeleteButtons();
    };
    return dsi;
}();
