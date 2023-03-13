"use strict";

function loadPosts(e) {
  document.getElementById("LoadMorePosts").disabled = true;
  fetch(window.location.href.split("?")[0] + "/loadPosts?after="+e).then(function (e) {
  if (e.ok) {
    return e.text();
  }
  }).then(function (e) {
    document.getElementById("posts").insertAdjacentHTML("beforeend", e);
    document.getElementById("LoadMorePosts").remove();  
  }).finally(function() {
    document.getElementById("LoadMorePosts").disabled = false;
  });
}