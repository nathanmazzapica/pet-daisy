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

/*
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
}*/

function buildMessage(name, content) {
    const message = document.createElement("p");
    message.classList.add("message");

    const sender = document.createElement("span");
    sender.innerText = `${name}: `;
    sender.classList.add("name");

    if (name !== displayName) {
        sender.classList.add("other");
    }

    const messageContent = document.createElement("span");

    const escapedContent = content.replace(/</g, "&lt;").replace(/>/g, "&gt;");

    const fragment = document.createDocumentFragment();
    const words = escapedContent.split(/\s+/);

    words.forEach(word => {
        if (/^https?:\/\/[^\s]+$/.test(word)) {
            const link = document.createElement("a");
            link.href = word.toLocaleLowerCase();
            link.target = "_blank";
            link.rel = "noopener noreferrer";
            link.textContent = word.toLocaleLowerCase();
            fragment.appendChild(link);
        } else {
            fragment.appendChild(document.createTextNode(word + " "));
        }
    });

    messageContent.appendChild(fragment);
    message.appendChild(sender);
    message.appendChild(messageContent);

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
        data: message
    }

    ws.send(JSON.stringify(chatMessage));
}