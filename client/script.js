const input = document.querySelector("#input");
const output = document.querySelector("#output");
const socket = new WebSocket("ws://138.68.58.77:2376/ws");

output.innerHTML += "Connecting...\n";

socket.onopen = () => {
    output.innerHTML
        += "Connected\n\n"
        + "Enter name\n";
};

socket.onerror = () => {
    output.innerHTML += "Failed\n";
};

socket.onmessage = (e) => {
    output.innerHTML += e.data;
};

input.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
        socket.send(input.value);
        input.value = "";
        input.focus();
        window.scrollTo(0, window.innerHeight)
    }
});

window.addEventListener("click", () => {
    input.focus();
});
