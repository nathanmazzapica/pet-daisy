const daisy = document.getElementById("daisy-image");
const counter = document.getElementById("counter");
const personalCounter = document.getElementById("personal-counter");
const chatInput = document.getElementById("chat-input");
const chatMessageContainer = document.getElementById("chat-message-container");

const ws = new WebSocket("ws://localhost:8080/ws")

let personalPets = 0;

daisy.addEventListener("click", () => {
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

window.addEventListener("resize", setGradientPosition);
window.addEventListener("load", setGradientPosition);