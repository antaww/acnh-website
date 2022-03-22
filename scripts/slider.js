var imgs = document.querySelectorAll('.slider-1 .sliderContent .houseware_box_item');
var prev = document.querySelector('.prev')
var next = document.querySelector('.next')
var sliderWraper = document.querySelector('.slider-1 .sliderContent')

var pos = 0;
var nbSlide = imgs.length;

window.addEventListener("keydown", event => {
    if (event.key == "ArrowLeft") {
        previousFunc()
    } else if (event.key == "ArrowRight") {
        nextFunc()
    }
})

prev.addEventListener('click', () => {
    previousFunc()
})

next.addEventListener('click', () => {
    nextFunc()
})

function previousFunc() {
    pos = pos + 100;
    if (pos > 0) {
        sliderWraper.style.left = ((-nbSlide + 1) * 100) + 'vw';
        pos = (-nbSlide + 1) * 100;
    } else {
        sliderWraper.style.left = pos + 'vw';
    }
}

function nextFunc() {
    pos = pos - 100;
    if (pos < ((-nbSlide + 1) * 100)) {
        sliderWraper.style.left = 0 + 'vw';
        pos = 0;
    } else {
        sliderWraper.style.left = pos + 'vw';
    }
}