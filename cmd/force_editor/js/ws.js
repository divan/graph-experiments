var ws = new WebSocket("ws://localhost:2018/ws", "graph");
ws.onopen = function (event) {
  exampleSocket.send("Here's some text that the server is urgently awaiting!"); 
};
