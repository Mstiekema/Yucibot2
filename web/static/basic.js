$(document).ready(function(){
  $("#searchBar").keyup(function(ev) {
    if (ev.which === 13) {
      var name = document.getElementById('searchBar').value
      location.href = "/user/" + name;
    }
  });
});