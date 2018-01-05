window.onload = function () {
  var conn = new ReconnectingWebSocket("ws://"+window.location.href.split("/")[2]+"/post/getCLR/", null, {debug: true, reconnectInterval: 3000});
  conn.onopen = function(e) {
    console.log("Connected forsenE")
  }
  conn.onclose = function(){
    console.log("xD")
  };
  conn.onmessage = function (e) {
    console.log(e)
    var data = JSON.parse(e.data)
    if (data.type == "meme") {
      var src = document.getElementsByTagName('source')[0]
      var video = document.getElementsByTagName('video')[0]
      src.setAttribute("src", data.meme)
      video.load();
      video.play();
      video.volume = 0.2;
      $(video).fadeIn(100);
      document.getElementsByTagName('video')[0].addEventListener('ended', function() {$("#meme").fadeOut(500);})
    } else if (data.type == "message") {
      console.log("Showed message: " + data.message)
      var id = "id" + (String(Math.floor(Math.random() * 1000000)))
      $("body").append("<p id="+id+">" + data.user + ": " + data.message + " </p><br>")
      $("#" + id).hide()
      $("#" + id).fadeIn(1000)
      document.getElementById(id).style.position = "absolute"
      document.getElementById(id).style.left = 100
      document.getElementById(id).style.top = 100
      setTimeout(function () { $("#" + id).fadeOut(1000); }, 5000);
      setTimeout(function () { document.getElementById(id).remove() }, 6000);
    } else if (data.type == "emote") {
      console.log("Showed emote: " + data.emote)
      var id = "id" + (String(Math.floor(Math.random() * 1000000)))
      $("body").append("<img id="+id+" src="+data.url+"><br>")
      $("#" + id).hide()
      var posx = (Math.random() * ($(document).width() - 300)).toFixed();
      var posy = (Math.random() * ($(document).height() - 300)).toFixed();
      $("#" + id).css({
        'position':'absolute',
        'left':posx+'px',
        'top':posy+'px',
        'display':'none'
      }).appendTo('body').fadeIn(1000)
      setTimeout(function () { $("#" + id).fadeOut(1000); }, 5000);
      setTimeout(function () { document.getElementById(id).remove() }, 6000);
    } else if (data.type == "sound") {
      console.log("Played sound: " + data.sound)
      var sound = new Audio(data.url);
      sound.play();
      sound.volume = 0.5;
    } else if (data.type == "gif") {
      console.log("Played gif: " + data.gif)
      var id = "id" + (String(Math.floor(Math.random() * 1000000)))
      $(".gif").append("<img style='width='300'; height='300';' id="+id+" src="+data.url+" align='middle'><br>")
      $("#" + id).hide()
      $("#" + id).fadeIn(1000)
      setTimeout(function () { $("#" + id).fadeOut(1000); }, 7000);
      setTimeout(function () { document.getElementById(id).remove() }, 8000);
    } else if (data.type == "chatEmote") {
      var id = "id" + (String(Math.floor(Math.random() * 1000000)))
      function getX() { return Math.floor(Math.random() * ($(window).width() - 300)) }
      function getY() { return Math.floor(Math.random() * $(window).height()) }
      $("body").append("<img id="+id+" src="+data.emote+"><br>")
      $("#" + id).hide()
      $("#" + id).fadeIn(1000)
      document.getElementById(id).style.position = "absolute"
      document.getElementById(id).style.left = getX()
      document.getElementById(id).style.top = getY()
      setTimeout(function () { $("#" + id).fadeOut(1000); }, 5000);
      setTimeout(function () { document.getElementById(id).remove() }, 6000);
    }
  }
}