/* main.js - handle events on the frontend */
function submitDate() {  // Called when user clicks the "Let's Go!" button
    console.log("submitted.");
    var month = document.getElementById("monthInput").value;
    var day = document.getElementById("dayInput").value;
    var start = document.getElementById("startInput").value;
    var end = document.getElementById("endInput").value;
    
    // Redirect to the get page with gym times, encoding the data in the URL
    window.location.href = "/get?month=" + month + "&day=" + day + "&start=" + start + "&end=" + end;   
};

function submitEmail() { 
    var email = document.getElementById("email").value;
    var url = new URLSearchParams(window.location.search);
    url.set('email', email);
    window.location.search = url;
}