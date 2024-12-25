const daisyContainer = document.getElementById('daisy-container');
const daisy = document.getElementById("daisy-image");
const counter = document.getElementById("counter");
const personalCounter = document.getElementById("personal-counter");
const chatInput = document.getElementById("chat-input");
const chatMessageContainer = document.getElementById("chat-message-container");

const ws = new WebSocket("ws://localhost:8080/ws")

let personalPets = 0;

let daisyReferenceSize = {width: 894, height: 597};
let referenceNoseCoordinates = {x: 349, y: 145};

const noseCoordinates = {x: 349, y: 145};
const noseHonk = new Audio("../static/honk.mp3");

const leftEyeCoordinate = {x: 320, y: 215};
const rightEyeCoordinate = {x: 505, y: 157};
const eyeOuch = new Audio("../static/ouch.mp3");

daisy.addEventListener("click", (e) => {

    const mousePos = {
        x: e.offsetX,
        y: e.offsetY
    }


    checkAndPerformEasterEggs(mousePos);


    petDaisy();
})

chatInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
        sendMessage(chatInput.value);
        chatInput.value = "";
    }
})

ws.onopen = () => {
    console.log("Connected to the server!");
};


ws.onmessage = (event) => {
    const data = JSON.parse(event.data);

    if (data.name == "petCounter") {
        counter.textContent = `Daisy has been pet ${data.message} times!`
        return;
    }

    if (data.name == "playerCount") {
        console.log("handle player count!")
        return;
    }

    if (data.name == "server") {
        console.log("handle server notification")
        console.log(data.message)
        chatMessageContainer.appendChild(displayServerChatNotification(data.message));
        return;
    }

    console.log(`${data.name}: ${data.message}`);

    chatMessageContainer.appendChild(buildMessage(data.name, data.message));
    chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight;

}

function petDaisy() {
    personalPets++;
    petMessage = {
        name: displayName,
        message: `$!pet;${personalPets}`,
    }
    ws.send(JSON.stringify(petMessage));
    personalCounter.innerText = `You have pet her ${personalPets} time${personalPets === 1 ? "" : "s"}!`;
}

function checkAndPerformEasterEggs(mousePos) {
    console.log(mousePos);
    if (inRadius(mousePos, noseCoordinates, 20)) {
        console.log("honk");
        noseHonk.currentTime = 0;
        noseHonk.play();
    }

    if (inRadius(mousePos, leftEyeCoordinate, 25) || inRadius(mousePos, rightEyeCoordinate, 5)) {
        console.log("ouch")
        eyeOuch.currentTime = 0;
        eyeOuch.play();
    }
}

function displayServerChatNotification(content) {
    const notification = document.createElement("p");
    notification.classList.add("notification", "message");

    if (content.indexOf(":(") !== -1) {
        notification.classList.add("notification", "disconnect");
    }

    notification.textContent = content.toUpperCase();
    return notification;
}

function buildMessage(name, content) {
    const message = document.createElement("p");
    message.classList.add("message");

    const spanClass = `${name !== displayName ? 'name' + ' other' : 'name'}`

    message.innerHTML = `<span class="${spanClass}">${name}:</span> ${content}`;

    return message;
}

function sendMessage(message) {

    message = message.trim();

    if (message === null || message.length === 0) {
        return;
    }

    chatMessage = {
        name: displayName,
        message
    }

    ws.send(JSON.stringify(chatMessage));
}

function setGradientPosition() {
    const daisyRect = daisy.getBoundingClientRect();

    const centerX = daisyRect.left + daisyRect.width / 2;
    const centerY = daisyRect.top + daisyRect.height / 2;

    document.body.style.background = `
        radial-gradient(circle at ${centerX}px ${centerY}px, #d0f0c0 1%, #0e250f)
    `;
}

function inRadius(point1, point2, radius) {
    const dist = Math.sqrt(Math.pow((point2.x - point1.x), 2) + Math.pow((point2.y - point1.y), 2));

    return dist <= radius;
}

function debugCircle(coordinate, id) {
    const circle = document.createElement("div");

    const rect = daisyContainer.getBoundingClientRect();
    let x = coordinate.x - rect.left;
    let y = coordinate.y - rect.top;

    circle.style.width = '20px';
    circle.style.height = '20px';
    circle.style.backgroundColor = 'red';
    circle.style.borderRadius = '50%';
    circle.style.position = 'absolute';
    circle.style.left = `${coordinate.x}px`;
    circle.style.top = `${coordinate.y}px`;
    circle.style.zIndex = '99999';
    circle.style.transform = 'translate(-50%, -50%)';

    if (id) {
        circle.id = id
    }


    daisy.appendChild(circle);

}

window.addEventListener("resize", () => {
    setGradientPosition();

    let daisySize = {width: daisy.width * 0.87, height: daisy.height * 0.87};
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
        setGradientPosition();

        let daisySize = {width: daisy.width * 0.87, height: daisy.height * 0.87};

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