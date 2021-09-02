let url = new URL('/ws',window.location.href)

if (url.protocol == 'http:') {
    url.protocol = url.protocol.replace('http','wss')
}
else
{
    url.protocol = url.protocol.replace('https','wss')
}
let socket = new WebSocket(url);
console.log("attempting connection")

socket.onopen = () => {
    console.log("Connected")
    socket.send("Hi from client")
}

socket.onclose = (event) => {
    console.log("Socket closed conn: ", event)
}

socket.onerror = (error) => {
    console.log("Socket error: ", error)
}

socket.onmessage = (msg) => {

}


$('#msbo').click( function (e){
    e.preventDefault();
    $(this).children("i").toggleClass('fa-bars').toggleClass('fa-times');
    $('body').toggleClass('msb-x');
});

function submitSetupForm(form){
    form.submit()
}
