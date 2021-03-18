/* main.js - handle events on the frontend */

function monthSelected() {
    console.log("YOU PICKED A MONTH.");
};

function submitDate() {  // Called when user clicks the "Let's Go!" button
    console.log("submitted.");
    var month = document.getElementById("monthSelector");
    var day = document.getElementById("dayInput").value;
    day = parseInt(day);
    var alertString;
    if (day && day > 0 && day <= 31) {
        alertString = "You chose: " + day;
    }
    else {
        alertString = "Invalid day selected!  Please enter a number in the text box and try again."
    }
    console.log(month.value);  //FIXME: not sure how to get the value from the dropdown yet
    console.log(day.value);

    // Write date to alert for now
    document.getElementById("alert-text").innerHTML = alertString;
};