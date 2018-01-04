$(document).ready(function(){
  $("#searchBar").keyup(function(ev) {
    if (ev.which === 13) {
      var name = document.getElementById('searchBar').value
      location.href = "/user/" + name;
    }
  });
  
  $("#fullProfile").click(function() {
    $(".profile").toggleClass("hidden")
    $("#name").toggleClass("hover")
  });
  
  document.onclick = function(event) {
    if (event.target.id == "name" || $(event.target).attr('class') == "profile text-left" || event.target.id == "menu" || $(event.target).attr('class') == "fa fa-caret-down") {
      return
    } else if( $(event.target).attr('class') == "img-circle pf") {
      window.location.href = "/user/" + event.target.id
    } else {
      $('.profile').addClass('hidden');
      $("#name").removeClass("hover")
    }
  }  
});