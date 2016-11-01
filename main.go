package main
import (
    "golang.org/x/net/websocket"
    "fmt"
    "log"
    "io"
    "net/http"
    "communicate/util"
)

var s *util.SockCon

func main() {
	/////////////////////
	s, _ = util.NewSock()
	go func() {
		cmd := "tail"
        	params := []string{"-f", "/var/log/messages"}	
		s.ExecCmd(cmd, params) 
	}()
    fmt.Println("Start app")
    http.Handle("/", websocket.Handler(Echo))
    http.HandleFunc("/client", Client)
    if err := http.ListenAndServe(":4001", nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func Echo(ws *websocket.Conn) {
    fmt.Println(ws)
    var err error
    bFlag := false
    for {
	log.Println("log")
        var recv string
        if err = websocket.Message.Receive(ws, &recv); err != nil {
            if err.Error() == "EOF" {
	    	fmt.Println("server receive EOF, client close the connect.")
		s.Mut.Lock()
	    	delete(s.MapStr, "jion")
		s.Mut.Unlock()
		ws.Close()
	    }else{
		fmt.Println("Can't receive", err)
	    }
            break
        }
        fmt.Println("Received from client: " + recv)
	
	if bFlag == false && recv == "start" {
        	fmt.Println("start web socket:")
		s.Mut.Lock()
		// add socketconnect
		s.MapStr["jion"] = ws
		if len(s.Content) != 0 {
			if err = websocket.Message.Send(ws, s.Content); err != nil {
                        	fmt.Println(err)
            }
		}
		s.Mut.Unlock()
		fmt.Println("s add jion and socket.")
		bFlag = true
	}
    }
}

func Client(w http.ResponseWriter, r *http.Request) {
	//html := `<html> <head> </head> <body> <script type="text/javascript"> var sock = null; var wsuri = "ws://192.168.56.101:4001"; window.onload = function() { console.log("onload"); sock = new WebSocket(wsuri); sock.onopen = function() { console.log("connected to " + wsuri); } sock.onclose = function(e) { console.log("connection closed (" + e.code + ")"); } sock.onmessage = function(e) { console.log("message received: " + e.data); } }; function send() { var msg = document.getElementById('message').value; alert("jion"); sock.send(msg); } </script> <h1>WebSocket Echo Test</h1> <form> <p>Message: <input id="message" type="text" value="Hello, world!"></p> </form> <button onclick="send();">Send Message</button> </body> </html>`	
	html := `<html> <head> </head> <body> <script type="text/javascript"> var sock = null; var wsuri = "ws://192.168.56.101:4001"; window.onload = function() { 
// checkout the browser support
window.WebSocket = window.WebSocket || window.MozWebSocket;
if (!window.WebSocket) {
	alert("WebSocket not supported by this browser.");
	return;
};

var opendiv = document.getElementById("open");
var logdiv = document.getElementById("log");
var datadiv =  document.getElementById("showData");
var clodiv = document.getElementById("close");
logdiv.innerText = "log" + logdiv.innerText;
console.log("onload");
sock = new WebSocket(wsuri); 
sock.onopen = function(){
	sendMsg("open");
	opendiv.innerText = "websocket open: build connect." + opendiv.innerText;
};
sock.onmessage = function(e){
	datadiv.innerText = datadiv.innerText + e.data;
};
sock.onclose = function(e){
	clodiv.innerText = "close connect.";
};
window.onbeforeunload = function() {
	closeSocket();
};

}; 
function send() { 
	var msg = document.getElementById('message').value; 
	sock.send(msg); 
};

function sendMsg(msg){
	
	sock.send(msg);
};
function closeSocket(){
	sock.close();
};
function start(){
	sendMsg("start")
}

</script>
 <h1>WebSocket Echo Test</h1> <form> <p>Message: <input id="message" type="text" value="Hello, world!"></p> </form> 
<button onclick="send();">Send Message</button>
<button onclick="closeSocket();"> close Socket </button>
<button onclick="start();"> start Socket </button>
<textarea>jion</textarea>
<div id="open"></div><div id="log"></div><div id="showData"></div>
<div id="close"></div>
 </body> </html>`	
	//html := `<html> <head> </head> <body> <script type="text/javascript"> var sock = null; var wsuri = "ws://192.168.56.101:4001";  function send() { var msg = document.getElementById('message').value;  } </script> <h1>WebSocket Echo Test</h1> <form> <p>Message: <input id="message" type="text" value="Hello, world!"></p> </form> <button onclick="send();"> Send Message </button></body> </html>`	
	//html := `<html> <head> </head> <body> 
 //<h1>WebSocket Echo Test</h1> <button ;">Send Message</button> </body> </html>`	
	io.WriteString(w, html)
}
