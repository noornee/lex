"use strict";

function loadPosts() {
  var loadPostForm = document.getElementById("LoadPostForm");
  var loadPostbtn = document.getElementById("LoadPostButton");

  if (loadPostForm != null && loadPostbtn != null) {
    loadPostbtn.disabled = true;
  
    return fetch("/loadPosts", {
      method: "POST",
      body: new FormData(loadPostForm)
    }).then(function(resp) {
      if (resp.ok) {
        return resp.text().then(function(ihtml) {
          document.getElementById("posts").insertAdjacentHTML('beforeend', ihtml);
          loadPostForm.remove();
        })
      } else {
        loadPostbtn.disabled = false;
      }
    }).catch(function() {
      loadPostbtn.disabled = false;
    }).finally(function(){
      loadPostbtn.disabled = false;
    })
  }
}

function loadUComm() {
  var loadPostForm = document.getElementById("LoadPostForm");
  var loadPostbtn = document.getElementById("LoadPostButton");

  if (loadPostForm != null && loadPostbtn != null) {
    loadPostbtn.disabled = true;
  
    return fetch("/loadAccComments", {
      method: "POST",
      body: new FormData(loadPostForm)
    }).then(function(resp) {
      if (resp.ok) {
        return resp.text().then(function(ihtml) {
          document.getElementById("posts").insertAdjacentHTML('beforeend', ihtml);
          loadPostForm.remove();
        })
      } else {
        loadPostbtn.disabled = false;
      }
    }).catch(function() {
      loadPostbtn.disabled = false;
    }).finally(function(){
      loadPostbtn.disabled = false;
    })
  }
}