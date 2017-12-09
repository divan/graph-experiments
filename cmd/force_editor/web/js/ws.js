var ws = new WebSocket('ws://' + window.location.host + '/ws');

var prev = function(e) {
	e.preventDefault();
	ws.send('{"cmd": "prev"}');
}

var next = function(e) {
	e.preventDefault();
	ws.send('{"cmd": "next"}');
}
