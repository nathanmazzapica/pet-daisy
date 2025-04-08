// TODO: export to .env
const ws = new WebSocket(wsURL)

let daisyReferenceSize = {width: 894, height: 597};
let referenceNoseCoordinates = {x: 349, y: 145};

/** @type {coordinate} */
const noseCoordinates = {x: 349, y: 145};
const noseHonk = new Audio("../static/honk.mp3");
noseHonk.volume = 0.1;
let nextHonk = 0;

const clickSound = new Audio("../static/audio/click.mp3");
clickSound.volume = 0.1;
clickSound.preload = "auto"
clickSound.currentTime = 0;


daisy.addEventListener("click", (e) => {

    const mousePos = {
        x: e.offsetX,
        y: e.offsetY
    }

    clickSound.currentTime = 0;
    clickSound.play();

    checkAndPerformEasterEggs(mousePos);
    petDaisy();
})


ws.onopen = () => {
    console.log("Connected to the server!");
};


ws.onmessage = (event) => {
    const message = JSON.parse(event.data);

    switch (message.name) {
        case "petCounter":
            handlePetCountUpdate(Number(message.data));
            break;
        case "playerCount":
            handlePlayerCountUpdate(message.data);
            break;
        case "server":
            handleServerNotification(message.data)
            break;
        case "leaderboard":
            displayLeaderboard(JSON.parse(message.data));
            break;
        default:
            handleIncomingChat(message.name, message.data)
    }

    if (message.name === "updateDisplay") {
        buildMessage("Daisy", `${displayName} has changed their name to ${message.data}!`)
        displayName = message.data;
        return
    }



}

ws.onclose = () => {
    window.location.href = "/error"
}

function petDaisy() {
    personalNumber++;
    petMessage = {
        name: displayName,
        data: `$!pet`,
    }
    ws.send(JSON.stringify(petMessage));
    personalCounter.innerText = `You have pet her ${Number(personalNumber).toLocaleString()} time${personalNumber === 1 ? "" : "s"}!`;
}

function handlePetCountUpdate(petCount) {
    let prettyCount = petCount.toLocaleString()
    counter.textContent = `Daisy has been pet ${prettyCount} times!`
}

function handlePlayerCountUpdate(playerCount) {
    document.getElementById("player-count").innerText = `Online Players: ${playerCount}`;
}

function handleServerNotification(content) {
    if (content.indexOf("hi!") !== -1 || content.indexOf(":(") !== -1) {
        chatMessageContainer.appendChild(displayServerChatNotification(content));
    }
    displayToast("Notification", content, 2000);
}

function handleIncomingChat(sender, content) {
    chatMessageContainer.appendChild(buildMessage(sender, content));
    chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight;

}

/**
 * @param {coordinate} mousePos
 */

function checkAndPerformEasterEggs(mousePos) {
    console.log(mousePos);
    if (inRadius(mousePos, noseCoordinates, 20) && nextHonk <= Date.now()) {
        console.log("honk");
        noseHonk.currentTime = 0;
        noseHonk.play();
        nextHonk = Date.now() + 500;
    }
}


function inRadius(point1, point2, radius) {
    const dist = Math.sqrt(Math.pow((point2.x - point1.x), 2) + Math.pow((point2.y - point1.y), 2));

    return dist <= radius;
}


window.addEventListener("resize", () => {
    const SCALE_OFFSET = 0.87
    if (window.innerWidth > 0) {
        setGradientPosition();
    } else {
        document.body.style.background = "var(--daisy-gradient-end)";
    }

    let daisySize = {width: daisy.width * SCALE_OFFSET, height: daisy.height * SCALE_OFFSET};
    console.log(daisySize)

    let ratio = {
        width: daisySize.width / daisyReferenceSize.width,
        height: daisySize.height / daisyReferenceSize.height
    };

    console.log(ratio)

    noseCoordinates.x = referenceNoseCoordinates.x * ratio.width;
    noseCoordinates.y = referenceNoseCoordinates.y * ratio.height;

    console.log(noseCoordinates.x)
    console.log(noseCoordinates.y)
});

window.addEventListener("load", () => {
        const SCALE_OFFSET = 0.87

        if (window.innerWidth > 0) {
            setGradientPosition();
        } else {
            document.body.style.background = "var(--daisy-gradient-end)";
        }

        let daisySize = {width: daisy.width * SCALE_OFFSET, height: daisy.height * SCALE_OFFSET};

        let ratio = {
            width: daisySize.width / daisyReferenceSize.width,
            height: daisySize.height / daisyReferenceSize.height
        };

        console.log(ratio)

        referenceNoseCoordinates.x = referenceNoseCoordinates.x * ratio.width;
        referenceNoseCoordinates.y = referenceNoseCoordinates.y * ratio.height;
        noseCoordinates.x = referenceNoseCoordinates.x;
        noseCoordinates.y = referenceNoseCoordinates.y;

        console.log(noseCoordinates.x)
        console.log(noseCoordinates.y)

        daisyReferenceSize = daisySize;
    }
);