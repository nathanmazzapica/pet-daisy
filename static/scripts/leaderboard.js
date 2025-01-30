/**
 * @param {leaderboardData[]} data
 * */
function displayLeaderboard(data) {

    const leaderboard = document.getElementById('leaderboard');
    leaderboard.innerHTML = '';
    for (let i = 0; i < data.length; i++) {
        const row = document.createElement('tr');
        const posCol = document.createElement('td');
        const displayNameCol = document.createElement('td');
        const petsCol = document.createElement('td');

        posCol.innerText = `${data[i].position}`;
        displayNameCol.innerText = `${data[i].display_name}`;
        petsCol.innerText = `${data[i].pet_count}`;

        if (data[i].display_name === displayName) {
            row.style.backgroundColor = 'var(--dark-purple-bg)';
        }

        row.appendChild(posCol);
        row.appendChild(displayNameCol);
        row.appendChild(petsCol);

        leaderboard.appendChild(row);
    }
}