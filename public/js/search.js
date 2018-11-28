$(document).ready(function() {
  $("#globalSearch").keypress(function(e) {
    if (e.which == 13) {
      document.location.href = "/explore/works/?q=" + $("#globalSearch").val();
    }
  });
});
