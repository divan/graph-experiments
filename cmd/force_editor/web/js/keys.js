document.addEventListener("keydown", function(event) {
	console.log(event.which) // keep this comment for the case when we need new keys
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
	if (event.which == 67) { // c
		calc(event);
		return;
	};
});
