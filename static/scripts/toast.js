function displayToast(title, msg, expiry) {
    const toast = document.createElement('div');
    toast.classList.add('toast','hide');

    const header = document.createElement('h4');
    header.textContent = title;
    const content = document.createElement('p');
    content.textContent = msg;

    toast.append(header);
    toast.append(content);
    document.body.appendChild(toast);

    setTimeout(() => {toast.classList.remove('hide')}, 100)

    setTimeout(() => {destroyToast(toast)}, expiry);
}

function destroyToast(toast) {
    toast.classList.add('fade')
    setTimeout(() => {
        toast.remove()
    }, 600);
}

//displayToast("Achievment", "Nathan has pet daisy 10,000 times!", 3500);