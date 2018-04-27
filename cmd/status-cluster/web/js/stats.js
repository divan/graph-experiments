var stats = {};

var Chart = require('chart.js');
let serversLineEl = document.getElementById('servers.sparkline');
var config = {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            data: [],
            backgroundColor: [ 'rgba(255, 99, 132, 0.2)' ,]
        }, {
            data: [],
            backgroundColor: [ 'rgba(255, 206, 86, 0.2)',]
        }]
    },
    options: {
        legend: {
            display: false
        },
        scales: {
            xAxes: [{
                display: false,
                type: 'category',
                stacked: true,
            }],
            yAxes: [{
                ticks: {
                    beginAtZero:true
                }
            }]
        }
    }
};
var serversChart = new Chart(serversLineEl, config);

function update(s) {
    stats = s;

    let serversDiv = document.getElementById('stats.servers');
    let clientsDiv = document.getElementById('stats.clients');
    serversDiv.innerHTML = stats.ServersNum;
    clientsDiv.innerHTML = stats.ClientsNum;

    // graphs
    serversChart.data.labels = stats.Timestamps ? stats.Timestamps : [];
    serversChart.data.datasets[0].data = stats.ServersHist;
    serversChart.data.datasets[1].data = stats.ClientsHist;
    serversChart.update();

    let table = [];

    if (stats.Nodes !== null) {
        stats.Nodes
            .filter(x => !x.IsClient)
            .map(x => table.push({ id: x.ID.slice(0, 6), peers: x.PeersNum, clients: x.ClientsNum }));
    }
    let tableBody = document.getElementById('servers.table');
    let out = '';
    table.forEach((x) => {
        out += '<tr>';
        out += '<td>' + x.id + '</td>';
        out += '<td>' + x.peers + '</td>';
        out += '<td>' + x.clients + '</td>';
        out += '</tr>';
    });
    tableBody.innerHTML = out;

    let lastUpdatedEl = document.getElementById('lastUpdated');
    lastUpdatedEl.innerHTML = new Date(stats.LastUpdate).toLocaleTimeString();
}

function current() {
    return stats;
}

module.exports = { update, current };
