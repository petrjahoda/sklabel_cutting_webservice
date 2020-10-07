const time = new EventSource('/time');
const goBack = document.getElementById('goback');
const result = document.getElementById('result');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
let entryData = ""
sessionStorage.clear()
sessionStorage.setItem("WorkplaceCode", workplace.textContent)
sessionStorage.setItem("User", user.textContent)
sessionStorage.setItem("UserId", user.dataset.userid)
sessionStorage.setItem("DeviceId", workplace.dataset.deviceid)

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

window.addEventListener("keyup", function (event) {
    entryData = entryData.replace("Meta", "");
    entryData = entryData.replaceAll("Enter", "");
    entryData = entryData.replaceAll("Shift", "");
    if (event.key === "Enter" && entryData.length > 0) {
        console.log("DATA: " + entryData)
        checkOrder(entryData);
        entryData = "";
    } else {
        entryData += event.key;
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
                let data = {
                    Type: "order",
                    Code: "K108",
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
                        Order: barcode,
                        DeviceId: sessionStorage.getItem("DeviceId"),
                        UserId: sessionStorage.getItem("UserId")
                    };
                    fetch("/create_order", {
                        method: "POST",
                        body: JSON.stringify(data)
                    }).then((response) => {
                        console.log("Starting order in Zapsi response: " + response.statusText);
                        window.location.replace('/home');
                    }).catch((error) => {
                        console.error('Error:', error);
                    });
                }).catch((error) => {
                    console.error('Error:', error);
                });

            } else {
                result.textContent = "Načtený kód " + barcode + " neexistuje v systému.";
                setTimeout(() => result.textContent = "Načtěte čárový kód zakázky", 3000)
            }
        });
    }).catch((error) => {
        console.error('Error:', error);
        result.textContent = "Načtený kód " + barcode + " neexistuje v systému.";
        setTimeout(() => result.textContent = "Načtěte čárový kód zakázky", 3000)
    });
}

