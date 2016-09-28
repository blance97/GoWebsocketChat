var Checked
$(document).ready(function() {
    checkLogin()
    listUserRoom()

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
    if (localStorage.getItem("RoomName") === null) {
        localStorage.setItem("RoomName", "room1");
    }
    $('#CurrentRoom').html("CurrentRoom: " + localStorage.getItem("RoomName"))
    getOldMessage(localStorage.getItem("RoomName"))
    $('.modal-trigger').leanModal();
    $(".button-collapse").sideNav();
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
                Roomname: localStorage.getItem("RoomName"),
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

var ws = new WebSocket("ws://" + window.location.host + "/entry");
ws.onopen = function() {
    $("#ChatPanel").html("CONNECTED")
}
ws.onclose = function() {
    $("#ChatPanel").html("DISCONNECTED")
    ws.close()
};
/**
This function is screwy because onload it loads previous messages so i have ot push to array TODO Fix
*/
ws.onmessage = function(event) {
    $("#ChatPanel").html("CONNECTED")
    var obj = jQuery.parseJSON(data)
    if (obj.Roomname == localStorage.getItem("RoomName")) {
        printJSON(event.data)
    }
}

function printJSON(data) {
    var obj = jQuery.parseJSON(data)
    Username = obj.author
    Text = obj.body
    $("#Chatbox").append('<p><b>' + Username + '</b>' + "(" + '<b>' + obj.time + '</b>' + "): " + Text + '</p>')
}

function getOldMessage(Roomname) {
    var request = $.ajax({
        type: 'GET',
        url: '/getOldMessage/?RoomName=' + Roomname,
        async: false,
        success: function(data) {
            var obj = jQuery.parseJSON(data)
            if (obj == null) {
                return
            }
            for (i = 0; i < obj.length; i++) {
                Username = obj[i].author
                Time = obj[i].time
                Test = obj[i].body
                $("#Chatbox").append('<p><b>' + Username + '</b>' + "(" + '<b>' + Time + '</b>' + "): " + Test + '</p>')
            }
        }
    });
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

function showUserInfo() {
    $('.button-collapse').sideNav('show');
}

function listUserRoom() {
    $("#users").html("")
    $.ajax({
        type: 'GET',
        url: '/listUsersinRoom/?RoomName=' + localStorage.getItem("RoomName"),
        async: false,
        success: function(data) {
            var obj = jQuery.parseJSON(data)
            if (obj == null) {
                return
            }
            for (i = 0; i < obj.length; i++) {
                $("#users").append('<a href="#" onclick="getUserInfo(\'' + obj[i] + '\')" data-activates="slide-out" class="button-collapse collection-item">' + obj[i] + '</a>')
            }
        },
        error: function(data) {
            alert("Could not recieve data from server")
        }

    });
}

function getUserInfo(user) {
    console.log("obtaining info for user: " + user)
    $.ajax({
        type: 'GET',
        url: '/getUserInfo/?Username=' + user,
        async: true,
        success: function(data) {
        $('#slide-out > li').remove();
          var obj = jQuery.parseJSON(data)
          $("#slide-out").append("<li><h5>Username: <b>" + user + "</h5></li>")
          $("#slide-out").append("<li><h5>"+ "Date Created:<br><b> " + timeConverter(obj.DateCreated) + "</h5></li>")
        }
    });
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

function updateUserRoom(RoomName) {
    $.ajax({
        type: 'POST',
        url: '/updateUserRoom',
        async: false,
        data: JSON.stringify({
            RoomName: RoomName
        }),
        dataType: 'json',
        success: function(data) {
            console.log("Succeed to update user room")
        },
        error: function(data) {
            console.log("Failed to update user room")
        }
    });
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
                updateUserRoom(RoomName)
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
        updateUserRoom(RoomName)
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
                    $("#RoomChanger").append('<a href="javascript:changews(\'' + obj.Rooms[i] + '\');" class="collection-item" style="display:inline-block;width:97%">' + obj.Rooms[i] + " (Private)" + '</a>');
                } else {
                    $("#RoomChanger").append('<a href="javascript:changews(\'' + obj.Rooms[i] + '\');" class="collection-item" style="display:inline-block;width:97%">' + obj.Rooms[i] + '</a>');
                }
                $("#RoomChanger").append('<a href="javascript:deleteRoom(\'' + obj.Rooms[i] + '\')"><i class="material-icons">delete_forever</i></a><br>')
            }
        }
    });
}

function deleteRoom(Room) {
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
function timeConverter(UNIX_timestamp) {

    var a = new Date(UNIX_timestamp * 1000);
    var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    var year = a.getFullYear();
    var month = months[a.getMonth()];
    var date = a.getDate();
    var hour = a.getHours();
    var min = "0" + a.getMinutes();
    var sec = "0" + a.getSeconds();
    var time = date + ' ' + month + ' ' + year + ' ' + hour + ':' + min.substr(-2) + ':' + sec.substr(-2);
    return time;
}
