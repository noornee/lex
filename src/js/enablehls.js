"use strict";

window.addEventListener('load', function(){
    loadHLS();
})

function loadHLS() {
    if (Hls.isSupported()) {
        var videos = document.querySelectorAll("video:not([hlsinit])")
        videos.forEach(function(v) {
            var playel = function() {
                var source = v.querySelector("source");

                if (source != null && source.src.includes("HLSPlaylist.m3u8"))
                {
                    var hls = new Hls({autoStartLoad: false});

                    hls.on(Hls.Events.MANIFEST_PARSED, function() {
                        var lvls = hls.levels;
                        if (lvls.length > 0) {
                            hls.startLevel = lvls.length - 1;
                        }
                    });

                    hls.loadSource(source.src);
                    hls.attachMedia(v);
                    hls.startLoad();

                    hls.on(Hls.Events.FRAG_BUFFERED, function() {
                        v.play();
                    })
                }
                v.removeEventListener("play", playel)
            }
            v.addEventListener("play", playel);
            v.setAttribute("hlsinit", "");
        })
    }
}