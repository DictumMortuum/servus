const monthNames = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December"
];

function startTime() {
  var today = new Date();
  var h = today.getHours();
  var m = today.getMinutes();
  var M = today.getMonth();
  var d = today.getDate();
  m = checkTime(m);
  document.getElementById('clock').innerHTML = h + " " + m;
  document.getElementById('date').innerHTML = monthNames[M] + " " + d;
  var t = setTimeout(startTime, 500);
}

function checkTime(i) {
  if (i < 10) {
    i = "0" + i
  };
  return i;
}

startTime();
