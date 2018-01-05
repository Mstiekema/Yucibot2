var conn;
var player;
var vids = new Array;
vids.Length = 0
conn = initWS();

function initWS() {
  var socket = new WebSocket("ws://"+window.location.href.split("/")[2]+"/post/getSongs/")
  socket.onopen = function (e) {
    conn.send("refreshData");
  };
  socket.onmessage = function (e) {
    var msg = JSON.parse(e.data)
    if (msg["Msg"] == "pushSonglist") {
      vids = JSON.parse(e.data)
    } else if (msg["Msg"] == "prevSongInfo") {
      player.loadVideoById({
        'videoId': msg["Songid"][0],
        'startSeconds': 0,
        'endSeconds': 600,
        'suggestedQuality': 'large'
      });
    } else if (msg == "'nextSong'") {
      nextVideo();
      conn.send('refreshData');
    } else if (msg == "'confDelSong'") {
      conn.send('refreshData');
    // } else if (msg == "'setVolume'") {
    //   var volume = parseInt(vol.volume)
    //   player.setVolume(volume)
    // } else if (msg == "'getVolume'") {
    //   var vol = player.getVolume()
    //   conn.send('returnVolume|'+vol)
    } else {
      console.log(e.data)
    }
  }
  return socket;
}

var tag = document.createElement('script');
tag.src = "https://www.youtube.com/iframe_api";
var firstScriptTag = document.getElementsByTagName('script')[0];
firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

function onYouTubeIframeAPIReady() {
  setTimeout(function() {
    player = new YT.Player('player', {
      height: '360',
      width: '640',
      videoId: vids.Songid[0],
      events: {
        'onReady': onPlayerReady,
        'onStateChange': onPlayerStateChange
      }
    });
  }, 500);
}

function onPlayerReady(event) {
  event.target.playVideo();
}

function onPlayerStateChange(event) {
  if (event.data == YT.PlayerState.ENDED) {
    conn.send('endSong|'+vids.Songid[0])
  } else if (event.data == -1) {
    conn.send("refreshData");
    event.target.playVideo();
  }
}

function nextVideo() {
	player.loadVideoById({
		'videoId': vids.Songid[1],
		'startSeconds': 0,
		'endSeconds': 600,
		'suggestedQuality': 'large'
	});
}

function previousVideo() {
	conn.send('prevSong')
}

function skipVideo() {
	conn.send('endSong|'+vids.Songid[0])
}

var getSongName = setInterval(function () {
	$("#videoTitle").html("<b>Current song</b>: " + vids.CurrSongT + "<br>" + "<b>Requested by</b>: " + vids.CurrSongN);
}, 5000);

var yucibot = angular.module('yucibot', [], function($interpolateProvider) {
    $interpolateProvider.startSymbol('*-');
    $interpolateProvider.endSymbol('-*');
});

yucibot.controller('songQueue', function($scope, $http, $log, $interval) {
  $scope.getVideos = function() {
    var total = 0;
    for(count = 0; count < vids.Length.length; count++){
      if (vids.Length[count] < 600) {
        total += parseInt(vids.Length[count])
      } else {
        total += 600
      }
    }
    var h = Math.floor(total / 3600);
    var m = Math.floor(total % 3600 / 60);
    var s = Math.floor(total % 3600 % 60);
    $scope.length = (h > 0 ? h + ":" + (m < 10 ? "0" : "") : "") + m + ":" + (s < 10 ? "0" : "") + s
    if(vids.Length[0]) {
      $scope.Ids = vids.Songid
      $scope.Thumbs = vids.Thumb
      $scope.Titles = vids.Title
      $scope.Names = vids.Name
    }
  }
  $scope.deleteVideo = function(id) {
    conn.send('delSong|'+id)
  }
  $scope.getVideos()
  $interval($scope.getVideos, 1000);
})