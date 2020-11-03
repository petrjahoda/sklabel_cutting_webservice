const time = new EventSource('/time');
const result = document.getElementById('result');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
const idleEnd = document.getElementById("idle-end");

workplace.textContent = sessionStorage.getItem("WorkplaceCode")
user.textContent = sessionStorage.getItem("User")
result.textContent = "Probíhá přestávka: " + sessionStorage.getItem("Idle")

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

idleEnd.addEventListener("touchend", endIdle)

function endIdle() {
    console.log("clicked")
    let data = {
        Type: "idle",
        Code: "K119",
        WorkplaceCode: workplace.textContent,
        UserId: sessionStorage.getItem("UserId"),
        OrderBarcode: sessionStorage.getItem("Order"),
        IdleId: sessionStorage.getItem("IdleId"),
    };
    fetch("/save_code", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        console.log("Saving code to K2 response: " + response.statusText);
        let data = {
            OrderBarcode: sessionStorage.getItem("Order"),
            IdleId: sessionStorage.getItem("IdleId"),
            DeviceId: sessionStorage.getItem("DeviceId"),
            UserId: sessionStorage.getItem("UserId")
        };
        fetch("/end_idle", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            console.log("Ending idle in Zapsi response: " + response.statusText);
            window.location.replace('/home');
        }).catch((error) => {
            console.error('Error:', error);
        });

    }).catch((error) => {
        console.error('Error:', error);
    });
}