$(function() {
	function prepareXHR() {
		var xhr = new window.XMLHttpRequest();
		xhr.upload.addEventListener("progress", function(evt) {
			if (evt.lengthComputable) {
				$pb.attr('data-transitiongoal', Math.ceil(evt.loaded* 100/evt.total)).progressbar({display_text: 'fill'});
			} else {
				console.error(evt);
			}
		}, false);
		return xhr;
	}
	var $pb = $('.progress .progress-bar');

	$("#_submit").on("click",function() {
		var files = $("#_file")[0].files;
		if(files.length === 0){
			return;
		}

		var data = new FormData();
		data.append('fupload', files[0]);

		$.ajax({
			xhr: prepareXHR,
			type: 'POST',
			url: "/upload",
			processData: false,
			contentType: false,
			data: data,
			success: function(data) {
				console.dir(data);
			},
			error: function(jqXHR, textStatus, errorMessage) {
				console.error(errorMessage); // Optional
      }
		});
	});
});