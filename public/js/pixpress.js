//'use strict';

// var csrf;
// var suburl;

$(document).ready(function () {
    // csrf = $('meta[name=_csrf]').attr("content");
    // suburl = $('meta[name=_suburl]').attr("content");

    // Dropzone
	var $dropzone = $('#dropzone');

    if ($dropzone.length > 0) {
        // Disable auto discover for all elements:
		Dropzone.autoDiscover = false;

        var filenameDict = {};
        $dropzone.dropzone({
			paramName: "files",
            url: $dropzone.data('upload-url'),
            // headers: {"X-Csrf-Token": csrf},
            maxFiles: $dropzone.data('max-file'),
            maxFilesize: $dropzone.data('max-size'),
            acceptedFiles: ($dropzone.data('accepts') === '*/*') ? null : $dropzone.data('accepts'),
            addRemoveLinks: true,
            dictDefaultMessage: $dropzone.data('default-message'),
            dictInvalidFileType: $dropzone.data('invalid-input-type'),
            dictFileTooBig: $dropzone.data('file-too-big'),
            dictRemoveFile: $dropzone.data('remove-file'),
            init: function () {
                this.on("success", function (file, data) {
                    filenameDict[file.name] = data.uuid;
                    var input = $('<input id="' + data.uuid + '" name="files" type="hidden">').val(data.uuid);
                    $('.files').append(input);
                });
                this.on("removedfile", function (file) {
                    if (file.name in filenameDict) {
                        $('#' + filenameDict[file.name]).remove();
                    }
                    if ($dropzone.data('remove-url') && $dropzone.data('csrf')) {
                        $.post($dropzone.data('remove-url'), {
                            file: filenameDict[file.name],
                            _csrf: $dropzone.data('csrf')
                        });
                    }
                })
            }
		});
	}
});
