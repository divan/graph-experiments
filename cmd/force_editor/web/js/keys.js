document.addEventListener("keydown", function(event) {
	if (event.which == 70) { // right
		toggleForces();
		return;
	};
	if (event.which == 221) { // right
		next(event);
		return;
	};
	if (event.which == 219) { // left
		prev(event);
		return;
	};
});
