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

const syncCodeForm = document.getElementById('sync-code-form');

syncCodeForm.addEventListener("submit", async function(event) {
    event.preventDefault();

    const code = document.getElementById("sync-code-input").value.trim();
    if (code === "") {
        alert("Please enter a sync code.");
        return;
    }

    try {
        const response = await fetch("/sync", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ code })
        });

        if (response.ok) {
            const data = await response.json();
            if (data.refresh) {
                window.location.reload();
            }
        } else {
            console.error("Failed to sync code:", response.statusText);
            alert("Failed to sync code. Please try again.");
        }
    } catch (error) {
        console.error("Error:", error);
        alert("An error occurred while syncing. Please try again.");
    }
});

