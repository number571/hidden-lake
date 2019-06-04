function show() {  
    $.ajax({  
        type: "POST",
        url: "/api/chat/global/update",  
        cache: false,  
        success: function(obj) {
            console.log(obj);
            if (obj.Exists) {
                $("#plain").html($("#plain").val() + obj.Body);  
            }
        },
    }).done(function(obj){
        show();
    });  
}  
$(document).ready(function(){  
    show();
});  
