/**
 * @param {leaderboardData[]} data
 * */
function displayLeaderboard(data) {
    leaderboardData = data;

    const leaderboard = document.getElementById('leaderboard');
    leaderboard.innerHTML = '';

    for (let i = 0; i < leaderboardData.length; i++) {
        const row = document.createElement('tr');
        const posCol = document.createElement('td');
        const displayNameCol = document.createElement('td');
        const petsCol = document.createElement('td');

        posCol.innerText = `${Number(leaderboardData[i].position).toLocaleString()}`;
        displayNameCol.innerText = `${leaderboardData[i].display_name}`;
        petsCol.innerText = `${leaderboardData[i].pet_count}`;

        if (leaderboardData[i].display_name === displayName) {
            row.style.backgroundColor = 'var(--dark-purple-bg)';
        }

        row.appendChild(posCol);
        row.appendChild(displayNameCol);
        row.appendChild(petsCol);

        leaderboard.appendChild(row);
    }
}

function updateLeaderboardDelta(changes) {
    for (const change of changes) {
        if (change.position === 0) {
            leaderboardData = leaderboardData.filter(r => r.display_name !== change.display_name);
            continue;
        }

        const idx = leaderboardData.findIndex(r => r.display_name === change.display_name);
        if (idx !== -1) {
            leaderboardData[idx] = change;
        } else {
            leaderboardData.push(change);
        }
    }

    leaderboardData.sort((a, b) => a.position - b.position);
    if (leaderboardData.length > 10) {
        leaderboardData = leaderboardData.slice(0, 10);
    }

    displayLeaderboard(leaderboardData);
}