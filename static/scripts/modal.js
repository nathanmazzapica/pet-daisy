const modalContainer = document.getElementById('modal-container');
const modal = document.getElementById('modal');



function openModal() {
    modalContainer.classList.remove('hidden')
}

function closeModal() {
    modalContainer.classList.add('hidden')
}

// put somewhere better later

const chatToggle = document.getElementById('chat-toggle');
const leaderboardToggle = document.getElementById('leaderboard-toggle');

let lbEnabled = true;
let chatEnabled = true;


// TODO: Make good
function checkResize() {
    setGradientPosition();
    if (!lbEnabled && !chatEnabled) {
        daisyContainer.style.width = "45%";
        return;
    }

    daisyContainer.style.width = "30%";

    if (isMobileDevice()) {
        daisyContainer.style.width = "25%";
    }
}

leaderboardToggle.addEventListener('change', () => {
    const leaderboard = document.getElementById('leaderboard-container');
    leaderboard.classList.toggle('hidden');
    lbEnabled = !lbEnabled;
    checkResize();
});

chatToggle.addEventListener('change', () => {
    const chat = document.getElementById('chat-container')
    chat.classList.toggle('hidden')
    chatEnabled = !chatEnabled;
    checkResize();
})

function isMobileDevice() {
    return /Mobi|Android|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
        || window.matchMedia("(max-width: 768px)").matches;
}

if (isMobileDevice()) {
    const chat = document.getElementById('chat-container')
    chat.classList.toggle('hidden')
    chatEnabled = false;
    chatToggle.checked = false;

    const leaderboard = document.getElementById('leaderboard-container');
    leaderboard.classList.toggle('hidden');
    lbEnabled = false;
    leaderboardToggle.checked = false;


    checkResize();
}

async function syncUserCode(code) {
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
                window.location.reload();  // Refresh the page
            }
        } else {
            console.error("Failed to sync code:", response.statusText);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}