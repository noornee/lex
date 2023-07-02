"use strict";

window.addEventListener('load', loadHLS)

function loadHLS() {
    if (!Hls.isSupported()) {
        return;
    }

    var videos = document.querySelectorAll("video:not([hlsinit])")

    videos.forEach(function(v) {
        var source = v.querySelector("source");

        if (source != null && source.src.includes("HLSPlaylist.m3u8"))
        {
            var hls = new Hls({autoStartLoad: false});

            hls.on(Hls.Events.MANIFEST_PARSED, function() {
                hls.startLevel = hls.levels.length - 1
            });
            
            v.addEventListener("play", function playV(){
                hls.loadSource(source.src);
                hls.attachMedia(v);
                hls.startLoad();
                v.play();

                v.removeEventListener("play", playV);
            });
            
            v.setAttribute("hlsinit", "");
        }
    })
}