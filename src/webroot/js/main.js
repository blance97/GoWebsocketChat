var Checked
$(document).ready(function() {
    checkLogin()
    $('#RoomName').val("")
    $('#test6').prop('checked', false);
    $('#test6').change(
        function() {
            if ($(this).is(':checked')) {
                $('#PrivatePass').show()
            } else {
                $('#PrivatePass').hide()
            }
        });
    for (i = 0; i < array.length; i++) {
        printJSON(array[i])
    }
    if (localStorage.getItem("RoomName") == null) {
        localStorage.setItem("RoomName", "room1");
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

var ws = new WebSocket("ws://" + window.location.host + "/entry/" + localStorage.getItem("RoomName"));
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

function checkPrivateRoom(RoomName) {
    var hasPerm
    $.ajax({
        type: 'GET',
        url: '/RoomPerm/?RoomName=' + RoomName,
        async: false,
        success: function(data) {
            hasPerm = true
        },
        error: function(data) {
            hasPerm = false
        }

    });
    console.log("Hasperm=", hasPerm)
    return hasPerm
}

function changews(RoomName) {
    if (checkPrivateRoom(RoomName)) {
        var Pass = prompt("Please enter your Password");
        if (Pass === null) {
            return
        }
        $.ajax({
            type: 'POST',
            url: '/RoomPassCheck',
            async: true,
            data: JSON.stringify({
                RoomName: RoomName,
                RoomPass: Pass
            }),
            dataType: 'json',
            success: function(data) {
                console.log("change room")
                localStorage.setItem("RoomName", RoomName);
                console.log(RoomName)
                window.location.reload()
            },
            error: function(data) {
                alert("Wrong Password")
            }
        });
    } else {
        console.log("change room")
        localStorage.setItem("RoomName", RoomName);
        console.log(RoomName)
        window.location.reload()
        return
    }
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
                if (obj.Private[i]) {
                    $("#RoomChanger").append('<a href="javascript:changews(\'' + obj.Rooms[i] + '\');" class="collection-item" style="display: inline-block;"width:97%>' + obj.Rooms[i] + ' (Private)</a>');
                } else {
                    $("#RoomChanger").append('<a href="javascript:changews(\'' + obj.Rooms[i] + '\');" class="collection-item" style="display:inline-block;width:97%">' + obj.Rooms[i] + '</a>');
                }
                $("#RoomChanger").append('<a href="javascript:deleteRoom(\'' + obj.Rooms[i] + '\')"><i class="material-icons">delete_forever</i></a><br>')
            }
        }
    });
}

function deleteRoom(Room){
  $.ajax({
      type: 'GET',
      url: '/deleteRoom/?RoomName=' + Room,
      async: true,
      success: function(data) {
        getRooms()
      },
      error: function(data) {
          alert("Could not delete")
      }
  });
}

function CreateRoom() {
    var Owner = getUser();
    var RoomName = $('#RoomName').val();
    var RoomPassword = $('#RoomPass').val();
    var Private
    if ($('#PrivatePass').is(':visible')) {
        Private = "true"
        if (RoomPassword == "") {
            $('#result').html("Password cannot be empty")
            return
        }
    } else {
        Private = "false"
        RoomPassword = ""
    }
    $.ajax({
        type: 'POST',
        url: '/createRoom',
        async: true,
        data: JSON.stringify({
            Owner: Owner,
            Roomname: RoomName,
            Private: Private,
            RoomPass: RoomPassword
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
