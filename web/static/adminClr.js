var soundEl = document.querySelectorAll('.sound');
var button = document.querySelectorAll('button');
var popup = document.getElementById('showBox');
var close = document.getElementById('close');
var memes = document.querySelectorAll('.memes');

window.onclick = function(event) {
  if (event.target == popup) {
    $('#showStuff').removeClass('animated bounceInDown');
    $('#showStuff').addClass('animated bounceOutUp');
    setTimeout(function () {
      popup.style.display = "none";
      window.location.reload()
    }, 700);
  }
}

if (close) {
  close.addEventListener('click', function() {
    $('#showStuff').removeClass('animated bounceInDown');
    $('#showStuff').addClass('animated bounceOutUp');
    setTimeout(function () {
      popup.style.display = "none";
      window.location.reload()
    }, 700);
  })
}

if (button) {
  for (var x = 0; x < button.length; x++) {
    button[x].addEventListener('click', function() {
      var type = $(this).attr("class")
      if(this.id == "loginBtn") return
      if (type == "btn btn-danger memes") {
        $.ajax({
          url: '/post/adminClr/',
          data: JSON.stringify({
            clrType: "meme"
          }),
          type: 'POST',
          success: function (data) {
            console.log(data)
          },
          error: function (xhr, status, error) {
            console.log("Error: "+xhr.responseText)
            // $("#notifBar").css({"background": "#f2dede", "color": "#a94442"}).fadeIn("slow").empty().append(xhr.responseText);
            // setTimeout(function () { $("#notifBar").fadeOut("slow"); }, 5000);
          }
        });
        return
      } else if (type.startsWith("btn btn-danger forceMeme")) {
        $.ajax({
          url: '/post/adminClr/',
          data: JSON.stringify({
            clrType: "forceMeme",
            name: this.id
          }),
          type: 'POST',
          success: function (data) {
            console.log(data)
          },
          error: function (xhr, status, error) {
            console.log("Error: "+xhr.responseText)
            // $("#notifBar").css({"background": "#f2dede", "color": "#a94442"}).fadeIn("slow").empty().append(xhr.responseText);
            // setTimeout(function () { $("#notifBar").fadeOut("slow"); }, 5000);
          }
        });
        return
      } else if (type.startsWith("btn btn-danger forceSound")) {
        $.ajax({
          url: '/post/adminClr/',
          data: JSON.stringify({
            clrType: "forceSound",
            name: this.id
          }),
          type: 'POST',
          success: function (data) {
            console.log(data)
          },
          error: function (xhr, status, error) {
            console.log("Error: "+xhr.responseText)
            // $("#notifBar").css({"background": "#f2dede", "color": "#a94442"}).fadeIn("slow").empty().append(xhr.responseText);
            // setTimeout(function () { $("#notifBar").fadeOut("slow"); }, 5000);
          }
        });
        return
      } else if (type.indexOf("btn btn-warning rem") != -1) {
        var id = type.slice(19)
        var con = confirm("Are you sure you want to remove " + this.id + "?");
        if (con == true) {
          $.ajax({
            url: '/post/adminClr/',
            data: JSON.stringify({
              clrType: "removeCLR",
              name: id
            }),
            type: 'POST',
            success: function (data) {
              console.log(data)
              window.location.reload()
            },
            error: function (xhr, status, error) {
              console.log("Error: "+xhr.responseText)
              // $("#notifBar").css({"background": "#f2dede", "color": "#a94442"}).fadeIn("slow").empty().append(xhr.responseText);
              // setTimeout(function () { $("#notifBar").fadeOut("slow"); }, 5000);
            }
          });
        }
        return
      } else if (type.indexOf("submitNewSample") != -1) {
        $.ajax({
          url: '/post/adminClr/',
          data: JSON.stringify({
            clrType: "addCLR",
            name: document.getElementById('id').value,
            url: document.getElementById('url').value,
            type: document.getElementById('type').value
          }),
          type: 'POST',
          success: function (data) {
            console.log(data)
            window.location.reload()
          },
          error: function (xhr, status, error) {
            console.log("Error: "+xhr.responseText)
            // $("#notifBar").css({"background": "#f2dede", "color": "#a94442"}).fadeIn("slow").empty().append(xhr.responseText);
            // setTimeout(function () { $("#notifBar").fadeOut("slow"); }, 5000);
          }
        });
        return
      }
      $('#showStuff').removeClass('animated bounceOutUp');
      $('#showStuff').addClass('animated bounceInDown');
      popup.style.display = "block";
      if (type.indexOf("pure-button gif") != -1) {
        var type = type.slice(16)
        $("#gifSpot").addClass(type)
        $("#gifSpot").addClass("gif")
        $("#gifSpot").addClass(this.id)
        $("#gifSpot").html("<img style='margin='auto'; width='300'; height='300';' src="+this.id+" align='middle'><br>")
      } else if (type.indexOf("pure-button sound") != -1) {
        var type = type.slice(18)
        $("#gifSpot").addClass(type)
        $("#gifSpot").addClass("sound")
        $("#gifSpot").addClass(this.id)
        $("#gifSpot").html("<img style='margin='auto'; width='300'; height='300';' src='https://www.clammr.com/Content/images/webapp/animated-sound.gif' align='middle'><br>")
      } else if (type.indexOf("btn btn-success clrSamples") != -1) {
        var clrType = type.slice(27)
        if (clrType == "gif") {
          $("#gifSpot").addClass(this.id)
          $("#gifSpot").addClass("gif")
          $("#gifSpot").html("<img style='margin='auto'; width='300'; height='300';' src="+this.id+" align='middle'><br>")
        } else if (clrType == "sound") {
          $("#gifSpot").addClass(this.id)
          $("#gifSpot").addClass("sound")
          $("#gifSpot").html("<img style='margin='auto'; width='300'; height='300';' src='https://www.clammr.com/Content/images/webapp/animated-sound.gif' align='middle'><br>")
        } else if (clrType == "meme") {
          $("#gifSpot").addClass(this.id)
          $("#gifSpot").addClass("meme")
          $("#gifSpot").html("<video style='margin='auto'; width='300'; height='300';' id='meme' width='640' height='360' align='middle' autoplay='true'><source src='"+this.id+"'type='video/mp4'></video>")
        } else {
          return
        }
      }
    })
  }
}

if (soundEl) {
  for (var x = 0; x < soundEl.length; x++) {
    soundEl[x].addEventListener('click', function() {
      var sound = new Audio(this.id);
      sound.play();
    })
  }
}