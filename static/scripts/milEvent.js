
const popSound = new Audio("../static/audio/pop_01.mp3");
popSound.volume = 0.5
popSound.preload = "auto"
function activateMilestoneEvent() {
    setTimeout(function () {
        displayToast("PARTY TIME", "Daisy is so excited to have been pet over ONE MILLION times and has invited all her friends to come play!", 15 * 1000)
    }, 100)

    const images = [
        "../static/images/DAISY.png",
        "../static/images/daisy_goofy.png",
        "../static/images/henry_goofy.png",
        "../static/images/Subject.png",
        "../static/images/daisy2.png",
        "../static/images/daisy3.png",
        "../static/images/happydaisy.png",
        "../static/images/henry.png",
        "../static/images/henry2.png",
        "../static/images/peter.png",
    ];

    const createFallingImage = () => {

        let clickable = Math.random() < 0.2 ? true : false;

        const img = document.createElement("img");
        img.src = images[Math.floor(Math.random() * images.length)]; // Pick a random image
        img.style.position = "fixed";
        img.style.top = "-100px";
        img.style.left = `${Math.random() * window.innerWidth}px`;
        img.style.width = "50px"; // Adjust as needed
        img.style.opacity = "0.8";
        img.style.zIndex = clickable ? "1" : "-1";
        img.style.transform = `rotate(${Math.random() * 360}deg)`;
        img.style.cursor = "pointer"; // Make it clickable
        img.draggable = false;
        img.style.userSelect = "none";

        if (clickable) {
            img.style.filter = "drop-shadow(4px 1px rgba(0, 0, 0, 0.8)) contrast(200%)";
        }

        const scale = Math.random() + 0.4;
        const rotation = Math.random() * 360;


        const fallDuration = Math.random() * 8 + 5;
        const rotateAmount = (Math.random() - 0.5) * 180;

        const fallAnimation = img.animate([
            { transform: `scale(${scale}) translateY(0) rotate(${Math.random() * 360}deg)`, opacity: 1 },
            { transform: `scale(${scale}) translateY(${window.innerHeight + 100}px) rotate(${rotateAmount}deg)`, opacity: 0 }
        ], {
            duration: fallDuration * 1000,
            easing: "linear",
            iterations: 1
        });

        fallAnimation.onfinish = () => {
            img.remove();
        };

        img.addEventListener("click", (event) => {
            createParticles(event.clientX, event.clientY);
            img.style.pointerEvents = "none"; // Prevent multiple clicks

            popSound.currentTime = 0; // Reset to start
            popSound.play();

            // Get current transform state to preserve translation & rotation
            const currentTransform = getComputedStyle(img).transform;

            img.animate([
                { transform: currentTransform, opacity: 1 },
                { transform: `${currentTransform} scale(1.5)`, opacity: 0 } // Scale relative to current state
            ], {
                duration: 300, // Quick pop
                easing: "ease-out",
                iterations: 1
            }).finished.then(() => img.remove());
            petDaisy()
            petDaisy()
            petDaisy()
            petDaisy()

        });
        document.body.appendChild(img);
    };

    setInterval(createFallingImage, 200); // Adjust spawn rate as needed
}

// document.addEventListener("DOMContentLoaded", () => {
//     const images = [
//         "../static/images/DAISY.png",
//         "../static/images/daisy_goofy.png",
//         "../static/images/henry_goofy.png",
//         "../static/images/Subject.png"
//     ];
//
//     const createFallingImage = () => {
//         const img = document.createElement("img");
//         img.src = images[Math.floor(Math.random() * images.length)]; // Pick a random image
//         img.style.position = "fixed";
//         img.style.top = "-100px";
//         img.style.left = `${Math.random() * window.innerWidth}px`;
//         img.style.width = "50px"; // Adjust as needed
//         img.style.opacity = "0.8";
//         img.style.zIndex = Math.random() < 0.3 ? "-1" : "1";
//         img.style.transform = `rotate(${Math.random() * 360}deg)`;
//         img.style.cursor = "pointer"; // Make it clickable
//         img.draggable = false;
//
//
//         const fallDuration = Math.random() * 8 + 5;
//         const rotateAmount = (Math.random() - 0.5) * 180;
//
//         const fallAnimation = img.animate([
//             { transform: `translateY(0) rotate(${Math.random() * 360}deg)`, opacity: 1 },
//             { transform: `translateY(${window.innerHeight + 100}px) rotate(${rotateAmount}deg)`, opacity: 0 }
//         ], {
//             duration: fallDuration * 1000,
//             easing: "linear",
//             iterations: 1
//         });
//
//         fallAnimation.onfinish = () => {
//             img.remove();
//         };
//
//         img.addEventListener("click", (event) => {
//             createParticles(event.clientX, event.clientY);
//             img.style.pointerEvents = "none"; // Prevent multiple clicks
//
//             // Get current transform state to preserve translation & rotation
//             const currentTransform = getComputedStyle(img).transform;
//
//             img.animate([
//                 { transform: currentTransform, opacity: 1 },
//                 { transform: `${currentTransform} scale(1.5)`, opacity: 0 } // Scale relative to current state
//             ], {
//                 duration: 300, // Quick pop
//                 easing: "ease-out",
//                 iterations: 1
//             }).finished.then(() => img.remove());
//
//         });
//         document.body.appendChild(img);
//     };
//
//     setInterval(createFallingImage, 200); // Adjust spawn rate as needed
// });

const createParticles = (x, y) => {
    const numParticles = 8; // Adjust for more/less particles

    for (let i = 0; i < numParticles; i++) {
        const particle = document.createElement("div");
        particle.classList.add("particle");

        // Position particle at click location
        particle.style.left = `${x}px`;
        particle.style.top = `${y}px`;

        // Random size
        const size = Math.random() * 6 + 4;
        particle.style.width = `${size}px`;
        particle.style.height = `${size}px`;

        // Random movement
        const angle = Math.random() * 2 * Math.PI;
        const distance = Math.random() * 50 + 30;
        const finalX = x + Math.cos(angle) * distance;
        const finalY = y + Math.sin(angle) * distance;

        document.body.appendChild(particle);

        // Animate particle
        particle.animate([
            { transform: `translate(0, 0)`, opacity: 1 },
            { transform: `translate(${finalX - x}px, ${finalY - y}px) scale(0.5)`, opacity: 0 }
        ], {
            duration: 500 + Math.random() * 300, // Randomized duration
            easing: "ease-out",
            iterations: 1
        }).finished.then(() => particle.remove());
    }
};

