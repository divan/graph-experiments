document.addEventListener("keydown", function(event) {
  console.log("key", event.which);
	if (event.which == 39) { // right
		next(event);
	};
	if (event.which == 37) { // left
		prev(event);
	};
});
