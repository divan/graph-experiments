var ngraph = require("ngraph.graph");
var graph = ngraph();
var graphData = require("../data/data.js").graphData;
var save = require('ngraph.tobinary');

graphData.nodes.forEach(node => { graph.addNode(node["id"]); });
graphData.links.forEach(link => { graph.addLink(link.source, link.target); });

// save meta.json, labels.json, links.bin
save(graph, {
  outDir: 'data/'
});

// generate positions.bin
var createLayout = require('ngraph.offline.layout');
console.log("createLayout start...");
var layout = createLayout(graph, {
  iterations: 500, // Run `100` iterations only
  saveEach: 500 // Save each `10th` iteration
});
layout.run();
