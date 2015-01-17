$(function() {
	var $pb = $('.progress .progress-bar');
	var panel = $("#uploadPanel");
	var status = $("#uploadStatus");

	function prepareXHR() {
		var xhr = new window.XMLHttpRequest();
		status.text("uploading");
		panel.addClass("panel-primary");
		xhr.upload.addEventListener("progress", function(evt) {
			if (evt.lengthComputable) {
				$pb.attr('data-transitiongoal', Math.ceil(evt.loaded* 100/evt.total)).progressbar({display_text: 'fill'});
			} else {
				status.text(evt);
				console.error(evt);
			}
		}, false);
		return xhr;
	}

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
				panel
					.removeClass("panel-primary")
					.addClass("panel-success");
				$(".progress").removeClass("active")
				$(".progress-bar").addClass("progress-bar-success");
				status.text(data);
			},
			error: function(jqXHR, textStatus, errorMessage) {
				panel
					.removeClass("panel-primary")
					.addClass("panel-danger");
				$(".progress").removeClass("active");
				$(".progress-bar").addClass("progress-bar-danger");
				var err = textStatus+":"+errorMessage;
				status.text(err);
				console.error(err);
				console.dir(jqXHR);
			}
		});
	});
});