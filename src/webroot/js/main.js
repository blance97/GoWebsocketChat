$(document).ready(function() {
    checkLogin()
    for (i = 0; i < array.length; i++) {
        printJSON(array[i])
    }
    $('#CurrentRoom').html("CurrentRoom: " + localStorage.getItem("RoomName"))
    $('.modal-trigger').leanModal();
    $('#loggedinAs').html("Current User: " + getUser())
    console.log("User: " + getUser())
    $("#inputChat").keyup(function(event) {
        if (event.keyCode == 13) {
            $("#sendMsgBtn").click();
        }
    });

    $('#inputChat').val("");
    $('#sendMsgBtn').click(function(event) {
        var d = new Date();
        var n = d.toLocaleTimeString();

        data = JSON.stringify({
            author: getUser(),
            time: n,
            body: $('#inputChat').val()
        }), ws.send(data);
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
    var Username
    var request = $.ajax({
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
function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}
function changews(RoomName) {
    console.log("change room")
    localStorage.setItem("RoomName", RoomName);
    console.log(RoomName)
    window.location.reload()
}

function getRooms() {
    $("#RoomChanger").html("")
    $.ajax({
        type: 'GET',
        url: '/getRooms',
        async: true,
        success: function(data) {
            var obj = jQuery.parseJSON(data)
            for (i = 0; i < obj.Rooms.length; i++) {
                $("#RoomChanger").append('<a href="javascript:changews(\'' + obj.Rooms[i] + '\');" class="collection-item" >' + obj.Rooms[i] + '</a>');
            }
        }
    });
}

function CreateRoom() {
    $.ajax({
        type: 'POST',
        url: '/createRoom',
        async: true,
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
        type: "GET",
        url: "http://localhost/logout",
        //data: {AppName: "Proline", Properties:null, Object: ""}, // An object, not a string.
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function(data) {
            window.location = "home.html"
        }
    })
}
