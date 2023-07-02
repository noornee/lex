"use strict";

window.addEventListener('load', initNav)

function initNav() {
    var gallery = document.querySelectorAll(".gallery:not([gallerynav])")

    gallery.forEach(function(thisgallery) {
        var inps = thisgallery.querySelectorAll("input[type='radio']")
        thisgallery.addEventListener("keydown", function(event){
            switch (event.key) {
                case "ArrowLeft":
                    var inpch = thisgallery.querySelector("input[type='radio']:checked")
                    var indxchecked = Array.from(inps).indexOf(inpch)
                    
                    if ((indxchecked-1) >= 0) {
                        inps[indxchecked-1].click();
                    }
                break;
                case "ArrowRight":
                    var inpch = thisgallery.querySelector("input[type='radio']:checked")
                    var indxchecked = Array.from(inps).indexOf(inpch)
                    
                    if ((indxchecked+1) < inps.length) {
                        inps[indxchecked+1].click();
                    }
                break;
            }
        })

        thisgallery.setAttribute("gallerynav", "")
    })
}
