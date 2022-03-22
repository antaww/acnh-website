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
