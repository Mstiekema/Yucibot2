var restart = document.getElementById("restart")
var switches = document.getElementById('.moduleSwitch');
$("[name='checkbox']").bootstrapSwitch();

$('.moduleSwitch').on('switchChange.bootstrapSwitch', function (event, state) {
  console.log(event)
  if (state == true) {
    sendData({"action": "enable", "module": this.id})
  } else {
    sendData({"action": "disable", "module": this.id})
  }
});

if (restart) {
  restart.addEventListener("click", function() {
    sendData({"action": "restart"})
  })
}

// This part of the code doesn't work yet.
var modChange = document.getElementById("modChange")
if (modChange) {
  modChange.addEventListener('click', function() {
    var newValues = [];
    $("input").each(function(){
      var id = $(this)[0].id; if (id == "searchBar") return; var value = document.getElementById(id).value
      if (value == "on") value = $(this).prop('checked'); if (value == true) value = 1; if (value == false) value = 0;
    });
    sendData({"action": "update", "id": id, "value": value})
  }, false);
}

function sendData(data) {
  $.ajax({
    url: '/post/modules/',
    data: JSON.stringify(data),
    type: 'POST',
    success: function (d) {
      console.log(d)
      window.location.reload()
    },
    error: function (xhr, status, error) {
      console.log("Error: "+xhr.responseText)
    }
  });
}
