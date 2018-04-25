var stats = {}

function update(stats) {
    stats= stats;
    let table = [];
    stats.Nodes
        .filter(x => !x.IsClient)
        .map(x => table.push({ id: x.ID.slice(0,6), peers: x.PeersNum, clients: x.ClientsNum }));
    let serversDiv = document.getElementById('stats.servers');
    let clientsDiv = document.getElementById('stats.clients');
    serversDiv.innerHTML = JSON.stringify(stats.ServersNum);
    clientsDiv.innerHTML = JSON.stringify(stats.ClientsNum);

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
}

module.exports = { update };
