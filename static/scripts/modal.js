const modalContainer = document.getElementById('modal-container');
const modal = document.getElementById('modal');



function openModal() {
    modalContainer.classList.remove('hidden')
}

function closeModal() {
    modalContainer.classList.add('hidden')
}