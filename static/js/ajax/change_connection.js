function change_connection() {
    var value = $('#selected_connection').val();
    console.log(value);
    if (value == "Global") { value = "" }
    window.location.replace("/network/chat/" + value)
}