var stats = {}

function update(stats) {
    stats= stats;
    let statsDiv = document.getElementById('stats');
    statsDiv.innerHTML = JSON.stringify(stats);
}

module.exports = { update };
