$(document).ready(function() {
    $("#myForm").submit(function(e) {
        var postData = $(this).serializeArray();
        var formURL = $(this).attr("action");
        $.ajax({
            url: formURL,
            type: "POST",
            data: postData,
            success: function(data, textStatus, jqXHR) {
              window.location = "index.html"
            },
            error: function(jqXHR, textStatus, errorThrown) {
                //if fails
            }
        });
        e.preventDefault(); //STOP default action
      //  e.unbind(); //unbind. to stop multiple form submit.
    });
    $("#ajaxform").submit(); //Submit  the FORM
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
//callback handler for form submit

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

function checkLogin() {
    $.ajax({
        type: 'GET',
        url: '/checkSession',
        async: false,
        error: function() {
            window.location = "index.html"
        }
    });
}
