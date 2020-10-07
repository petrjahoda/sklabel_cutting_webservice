const time = new EventSource('/time');
const result = document.getElementById('result');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
let data = ""
workplace.textContent = sessionStorage.getItem("WorkplaceCode")
user.textContent = sessionStorage.getItem("User")
result.textContent = "Přihlaste se přiložením karty"

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

window.addEventListener("keyup", function (event) {
    data = data.replace("Meta", "");
    data = data.replaceAll("Enter", "");
    data = data.replaceAll("Shift", "");
    if (event.key === "Enter" && data.length > 0) {
        console.log("DATA: " + data)
        checkUser(data);
        data = "";
    } else {
        data += event.key;
    }
});

function checkUser(barcode) {
    let data = {Data: barcode};
    fetch("/check_user", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (text) {
            let myObj = JSON.parse(text);
            let currentResult = myObj.Data;
            console.log("0")
            if (currentResult === "ok") {
                let data = {
                    Order: barcode,
                    DeviceId: sessionStorage.getItem("DeviceId"),
                    UserId: sessionStorage.getItem("UserId")
                };
                fetch("/end_order", {
                    method: "POST",
                    body: JSON.stringify(data)
                }).then((response) => {
                    console.log("Ending order in Zapsi response: " + response.statusText);
                    let data = {
                        Type: "order",
                        Code: "K105",
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
                        sessionStorage.setItem("UserId", myObj.UserId)
                        sessionStorage.setItem("User", "Přihlášen: " + myObj.UserName)
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
                }).catch((error) => {
                    console.error('Error:', error);
                });

            } else {
                result.textContent = "Uživatel " + barcode + " neexistuje v systému";
                setTimeout(() => result.textContent = "Přihlaste se přiložením karty", 30000)
            }
        });
    }).catch((error) => {
        console.error('Error:', error);
        result.textContent = "Uživatel " + barcode + " neexistuje v systému";
        setTimeout(() => result.textContent = "Přihlaste se přiložením karty", 30000)
    });
}