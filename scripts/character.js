const house_interior = document.querySelector(".house_interior")
const dialogue = document.querySelector(".dialogue")
house_interior.addEventListener('click', event => {
    console.log("house interior clicked")
    house_interior.classList.toggle("house_interior_fullsize")
    dialogue.style.zIndex = "1";
})