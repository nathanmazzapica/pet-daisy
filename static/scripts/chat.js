const chatInput = document.getElementById("chat-input");
const chatMessageContainer = document.getElementById("chat-message-container");

chatInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
        sendMessage(chatInput.value);
        chatInput.value = "";
    }
})

function displayServerChatNotification(content) {
    const notification = document.createElement("p");
    notification.classList.add("notification", "message");

    if (content.indexOf(":(") !== -1) {
        notification.classList.add("notification", "disconnect");
    }

    if (content.indexOf("say hi!") !== -1) {
        notification.classList.add("notification", "connect")
    }

    notification.textContent = content.toUpperCase();
    return notification;
}

function buildMessage(name, content) {
    const message = document.createElement("p");
    message.classList.add("message");

    const sender = document.createElement("span")
    sender.innerText = `${name}: `;

    sender.classList.add('name');
    if (name !== displayName) {
        sender.classList.add('other')
    }

    message.textContent = content;
    message.prepend(sender);

    return message;
}

function sendMessage(message) {

    message = message.trim();

    if (message === null || message.length === 0) {
        return;
    }

    if (message.length > 256) {
        alert("Message is too long!");
        return;
    }

    chatMessage = {
        name: displayName,
        message
    }

    ws.send(JSON.stringify(chatMessage));
}