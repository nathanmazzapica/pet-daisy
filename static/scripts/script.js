// TODO: export to .env
const ws = new WebSocket(wsURL)

let daisyReferenceSize = {width: 894, height: 597};
let referenceNoseCoordinates = {x: 349, y: 145};

/** @type {coordinate} */
const noseCoordinates = {x: 349, y: 145};
const noseHonk = new Audio("../static/honk.mp3");
noseHonk.volume = 0.1;
let nextHonk = 0;


daisy.addEventListener("click", (e) => {

    const mousePos = {
        x: e.offsetX,
        y: e.offsetY
    }

    checkAndPerformEasterEggs(mousePos);
    petDaisy();
})


ws.onopen = () => {
    console.log("Connected to the server!");
};


ws.onmessage = (event) => {
    const data = JSON.parse(event.data);

    if (data.name === "petCounter") {
        let prettyCount = Number(data.message).toLocaleString()
        counter.textContent = `Daisy has been pet ${prettyCount} times!`
        return;
    }

    if (data.name === "playerCount") {
        console.log("handle player count!")
        console.log(data.message);
        document.getElementById("player-count").innerText = `Online Players: ${data.message}`;
        return;
    }

    if (data.name === "server") {
        console.log("handle server notification")
        console.log(data.message)
        displayToast("Notification", data.message, 2000);
        chatMessageContainer.appendChild(displayServerChatNotification(data.message));
        chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight;
        return;
    }

    if (data.name === "leaderboard") {
        console.log("handle leaderboard notification")
        console.log(JSON.parse(data.message));
        displayLeaderboard(JSON.parse(data.message));
        chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight;
        return;
    }

    console.log(`${data.name}: ${data.message}`);

    chatMessageContainer.appendChild(buildMessage(data.name, data.message));
    chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight;

}

ws.onclose = () => {
    window.location.href = "/error"
}

function petDaisy() {
    personalNumber++;
    petMessage = {
        name: displayName,
        message: `$!pet;${personalNumber}`,
    }
    ws.send(JSON.stringify(petMessage));
    personalCounter.innerText = `You have pet her ${personalNumber} time${personalNumber === 1 ? "" : "s"}!`;
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