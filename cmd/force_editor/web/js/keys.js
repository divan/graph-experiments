document.addEventListener("keydown", function(event) {
	if (event.which == 221) { // right
		next(event);
	};
	if (event.which == 219) { // left
		prev(event);
	};
});
