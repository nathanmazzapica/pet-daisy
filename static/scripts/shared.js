/**
 *  @typedef {{x: number, y: number}} coordinate
 */

/**
 *
 * @typedef {{display_name: string, pet_count: number, position: number}} leaderboardData
 */


const daisyContainer = document.getElementById('daisy-container');
const daisy = document.getElementById("daisy-image");
const counter = document.getElementById("counter");
const personalCounter = document.getElementById("personal-counter");


function setGradientPosition() {
    const daisyRect = daisy.getBoundingClientRect();

    const centerX = daisyRect.left + daisyRect.width / 2;
    const centerY = daisyRect.top + daisyRect.height / 2;

    document.body.style.background = `
        radial-gradient(circle at ${centerX}px ${centerY}px, var(--daisy-gradient-start) 1%, var(--daisy-gradient-end) 100%
    `;
}
