// const arrow = document.querySelector(".arrow")
//
// arrow.addEventListener("click", () => {
//         console.log("arrow clicked")
//     }
// )

window.addEventListener("load", function () {
    function isVisibleInViewport(element) {
        if (element.offsetWidth || element.offsetHeight || element.getClientRects().length) {
            const rect = element.getBoundingClientRect();
            return rect.bottom > 0 && rect.right > 0 && rect.left < (window.innerWidth || document.documentElement.clientWidth) && rect.top < (window.innerHeight || document.documentElement.clientHeight);
        }
        return false;
    }

    function animateVisibleElements() {
        const list = document.querySelectorAll('.chara_box');

        list.forEach(item => {
            if (isVisibleInViewport(item)) item.classList.add('animate');
        });
    }

    document.addEventListener('scroll', animateVisibleElements);
    animateVisibleElements();

});


const menu_btn = document.querySelector(".button")
menu_btn.addEventListener('click', event => {
    console.log("menu button clicked")
    document.body.classList.toggle("no-scroll")
    console.log("body noscroll")
})

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
