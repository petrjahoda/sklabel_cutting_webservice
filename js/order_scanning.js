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


document.addEventListener("keyup", function (event) {
    entryData += event.key;
    if (entryData.includes("Enter")) {
        checkOrder(entryData.toUpperCase());
        entryData = ""
    }
});


document.addEventListener("touchend", function (event) {
    result.textContent = "touched"
});

goBack.addEventListener("touchend", () => {
    window.history.back();
})
// goBack.addEventListener("click", () => {
//     window.history.back();
// })

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
                sessionStorage.setItem("Order", myObj.Result)
                let data = {
                    Type: "order",
                    Code: "108",
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
                        OrderBarcode: myObj.Result,
                        DeviceId: sessionStorage.getItem("DeviceId"),
                        UserId: sessionStorage.getItem("UserId"),
                        Pcs: "0",
                        CloseLogin: "false"
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
                result.textContent = "Načtený kód " + myObj.Result + " neexistuje v systému.";
                setTimeout(() => result.textContent = "Načtěte čárový kód zakázky", 3000)
            }
        });
    }).catch((error) => {
        console.error('Error:', error);
        result.textContent = "Chyba komunikace";
        setTimeout(() => result.textContent = "Načtěte čárový kód zakázky", 3000)
    });
}

window.focus()
