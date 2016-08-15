$(document).ready(function() {
    for (i = 0; i < array.length; i++) {
        printJSON(array[i])
    }
    $('#CurrentRoom').html("CurrentRoom: " + localStorage.getItem("RoomName"))
    $('.modal-trigger').leanModal();
    $('#loggedinAs').html("Current User: " + getUser())
    $('#inputChat').val("");
    $('#sendMsgBtn').click(function(event) {
        var d = new Date();
        var n = d.toLocaleTimeString();
        data = JSON.stringify({
                author: getUser(),
                time: n,
                body: $('#inputChat').val()
            }),
            ws.send(data);
        console.log($('#inputChat').val())
        $('#inputChat').val("");
        scrollBottom()
    });

});


function scrollBottom() {
    var height = 0;
    $('div p').each(function(i, value) {
        height += parseInt($(this).height());
    });

    height += '';
    $('div').animate({
        scrollTop: height
    });
}

ws = new WebSocket("ws://localhost:80/entry/" + localStorage.getItem("RoomName"));

function changews(RoomName) {
    console.log("change room")
    localStorage.setItem("RoomName", RoomName);
    console.log(RoomName)
    window.location.reload()
}

ws.onopen = function() {
    $("#ChatPanel").html("CONNECTED")
};
ws.onclose = function() {
    $("#ChatPanel").html("DISCONNECTED")
    ws.close()
};
var array = []
    /**
    This function is screwy because onload it loads previous messages so i have ot push to array TODO Fix
    */
ws.onmessage = function(event) {
  //  console.log("Recieved Message: " + event.data)
    array.push(event.data)
    printJSON(event.data)

}

function printJSON(data) {
    var obj = jQuery.parseJSON(data)
    Username = obj.author
    Text = obj.body
    $("#Chatbox").append('<p><b>' + Username + '</b>' + "(" + '<b>' + obj.time + '</b>' + "): " + Text + '</p>')
}

function getUser() {
    var result = null;
    $.ajax({
        type: 'GET',
        url: '/getUser',
        async: false,
        success: function(data) {
            var obj = jQuery.parseJSON(data)
            Username = obj.Username
        }
    });
    return Username
}

function CreateRoom() {
    $.ajax({
        type: 'POST',
        url: '/createRoom',
        data: JSON.stringify({
            RoomName: $('#RoomName').val()
        }),
        dataType: 'json',
        success: function(data) {
            console.log("Posted Data");
            $('#modal1').closeModal();
        }
    });
}

function logout() {
    $.ajax({
        type: 'GET',
        url: '/logout',
        async: false,
        success: function() {
            window.location = "index.html"
        }
    });
}
