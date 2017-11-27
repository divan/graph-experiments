var ngraph = require("ngraph.graph");
var graph = ngraph();
var graphData = require("./data.js").graphData;
console.log(graphData.nodes);

graphData.nodes.forEach(node => { graph.addNode(node["id"]); });
graphData.links.forEach(link => { graph.addLink(link.source, link.target); });

var createLayout = require('ngraph.offline.layout');
console.log("createLayout start...");
var layout = createLayout(graph, {
  iterations: 500, // Run `100` iterations only
  saveEach: 500 // Save each `10th` iteration
});
layout.run();
