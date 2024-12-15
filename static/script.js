const daisy = document.getElementById("daisy-image");
const counter = document.getElementById("counter");
const personalCounter = document.getElementById("personal-counter");
const chatInput = document.getElementById("chat-input");
const chatMessageContainer = document.getElementById("chat-message-container");

const ws = new WebSocket("ws://localhost:8080/ws")

let personalPets = 0;

const noseCoordinates = { x: 215, y: 90 };
const noseHonk = new Audio("../static/honk.mp3");

const leftEyeCoordinate = { x: 320, y: 215 };
const rightEyeCoordinate = { x: 505, y: 157 };
const eyeOuch = new Audio("../static/ouch.mp3");

daisy.addEventListener("click", (e) => {

    const rect = daisy.getBoundingClientRect();
    let x = e.clientX - rect.left;
    let y = e.clientY - rect.top;
    console.log(x, y);
    if (inRadius({x, y}, noseCoordinates, 13)) {
        console.log("honk");
        noseHonk.currentTime = 0;
        noseHonk.play();
    }

    if (inRadius({x, y}, leftEyeCoordinate, 25) || inRadius({x, y}, rightEyeCoordinate, 5)) {
        eyeOuch.currentTime = 0;
        eyeOuch.play();
    }

    // extract to func
    personalPets++;
    petMessage = {
        name: "nathan",
        message: `$!pet;${personalPets}`,
    }
    ws.send(JSON.stringify(petMessage));
    personalCounter.innerText = `You have pet her ${personalPets} time${personalPets === 1 ? "" : "s"}!`;
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

}

function displayServerChatNotification(content) {
    const notification = document.createElement("p");
    notification.classList.add("notification", "message");

    notification.textContent = content.toUpperCase();
    return notification;
}

function buildMessage(name, content) {
    const message = document.createElement("p");
    message.classList.add("message");

    message.innerHTML = `<span class="name">${name}:</span> ${content}`;

    return message;
}

function sendMessage(message) {

    message = message.trim();

    if (message === null || message.length === 0) {
        return;
    }

    chatMessage = {
        name: "nathan",
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
    const dist = Math.sqrt(Math.pow((point2.x - point1.x), 2) + Math.pow((point2.y - point1.y),2));

    return dist <= radius;
}

window.addEventListener("resize", setGradientPosition);
window.addEventListener("load", setGradientPosition);