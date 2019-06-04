$(document).ready(function(){  
    $('#mode').click(function(){  
        $.ajax({  
            type: "POST",  
            url: "/api/mode",
            cache: false,
            success: function(obj){
                if (obj.Exists) {
                    location.reload();
                }
            }  
        });
        return false;
    });  
});  
