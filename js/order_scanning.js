const time = new EventSource('/time');
const goBack = document.getElementById('goback');
const result = document.getElementById('result');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
let data = ""

sessionStorage.clear()
sessionStorage.setItem("WorkplaceCode", workplace.textContent)
sessionStorage.setItem("User", user.textContent)
sessionStorage.setItem("UserId", workplace.dataset.deviceid)
sessionStorage.setItem("DeviceId", user.dataset.userid)

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);


window.addEventListener("keyup", function (event) {
    data = data.replace("Meta", "");
    data = data.replaceAll("Enter", "");
    if (event.key === "Enter" && data.length > 0) {
        console.log("DATA: " + data)
        checkOrder(data);
        data = "";
    } else {
        data += event.key;
    }
});

goBack.addEventListener("touchend", () => {
    window.history.back();
})
goBack.addEventListener("click", () => {
    window.history.back();
})

function checkOrder(barcode) {
    let data = {Data: barcode};
    fetch("/check_order", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (text) {
            let myObj = JSON.parse(text);
            let currentResult = myObj.Data;
            if (currentResult === "ok") {
                sessionStorage.setItem("Order", barcode)
                saveCodeToK2("K108")
                startOrderInZapsi(barcode, workplace.dataset.deviceid, user.dataset.userid);
                window.location.replace('/home');
            } else {
                result.textContent = "Načtený kód " + barcode + " neexistuje v systému.";
            }
        });
    }).catch((error) => {
        console.error('Error:', error);
        result.textContent = "Načtený kód " + barcode + " neexistuje v systému.";
    });
}

function startOrderInZapsi(barcode, deviceid, userid) {
    let data = {Order: barcode, DeviceId: deviceid, UserId: userid};
    fetch("/start_order", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        console.log(response.statusText);
    }).catch((error) => {
        console.error('Error:', error);
    });
}


function saveCodeToK2(code) {
    let data = {Data: code};
    fetch("/save_code", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        console.log(response.statusText);
    }).catch((error) => {
        console.error('Error:', error);
    });
}