$(function () {
	$("#content").bind("input change", function () {
		$.post("/write", {md: $("#content").val()}, function (response) {
			$("#md_html").html(response.html)
		});
	});
})