$(document).ready(function() {
    // the "href" attribute of .modal-trigger must specify the modal ID that wants to be triggered
    $('.modal-trigger').leanModal();
    $("#myForm").bind('ajax:complete', function() {
        alert(hi)

    });

});

function validateForm() {

    if ($("#PasswordSignup").val() != $("#PasswordSignup2").val()) {
        console.log("passwords don't match")
        $("#result").html("Passwords Do Not Match")
    } else {
        console.log("passwords are good")
        $("#result").html("")
    }
}

function signup() {
    $.ajax({
        type: 'POST',
        url: '/signup',
        data: JSON.stringify({
            Username: $('#UsernameSignup').val(),
            Pass: $('#PasswordSignup').val()
        }),
        dataType: 'json',
        async: false,
        success: function(data) {
            console.log("Posted Data");
        },
        error: function(xhr, textStatus, error) {
            console.log(xhr.statusText);
            console.log(textStatus);
            console.log(error);
        }
    });
}

function login() {
    $.ajax({
        type: 'POST',
        url: '/login',
        data: JSON.stringify({
            Username: $('#Username').val(),
            Pass: $('#password').val()
        }),
        dataType: 'json',
        async: false,
        success: function(data) {
            console.log("Posted Data");
        },
        error: function(xhr, textStatus, error) {
            console.log(xhr.statusText);
            console.log(textStatus);
            console.log(error);
        }
    });

}

function checkLogin(){
      $.ajax({
          type: 'GET',
          url: '/checkSession',
          async: false,
          error: function() {
              window.location = "index.html"
          }
      });
}
