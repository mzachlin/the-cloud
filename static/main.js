/* main.js - handle events on the frontend */

function monthSelected() {
    console.log("YOU PICKED A MONTH.");
};

function submitDate() {  // Called when user clicks the "Let's Go!" button
    console.log("submitted.");
    var month = document.getElementById("monthInput").value;
    var day = document.getElementById("dayInput").value;
    month = parseInt(month);
    day = parseInt(day);
    var alertString;
    if (day && day > 0 && day <= 31) {
        if (month && month > 0 && month <= 12) {
            alertString = "You chose month " + month + " and day " + day;
            // Redirect to the get page with gym times, encoding the data in the URL
            console.log(month);
            console.log("testing ya");
            window.location.href = "/get?month=" + month + "&day=" + day;
        }
        else {
            alertString = "Invalid month selected!  Please enter a number in the text box and try again."
        }
        
    }
    else {
        alertString = "Invalid day selected!  Please enter a number in the text box and try again."
    }
    //console.log(month.value);  //FIXME: not sure how to get the value from the dropdown yet
    console.log(day);

    // Write date to alert for now
    document.getElementById("alert-text").innerHTML = alertString;
};