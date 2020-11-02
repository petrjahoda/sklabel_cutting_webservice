let data = {
    OrderBarcode: sessionStorage.getItem("Order"),
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
        Code: "K219",
        WorkplaceCode: sessionStorage.getItem("WorkplaceCode"),
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
            Type: "order",
            Code: "0004",
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
            sessionStorage.setItem("UserId", "")
            sessionStorage.setItem("User", "Přihlášen: " + "")
            user.textContent = "Nenávaznost obsluhy"
            let data = {
                OrderBarcode: sessionStorage.getItem("Order"),
                DeviceId: sessionStorage.getItem("DeviceId"),
                UserId: sessionStorage.getItem("")
            };
            fetch("/create_order", {
                method: "POST",
                body: JSON.stringify(data)
            }).then((response) => {
                console.log("Starting order in Zapsi response: " + response.statusText);
            }).catch((error) => {
                console.error('Error:', error);
            });

        }).catch((error) => {
            console.error('Error:', error);
        });

    }).catch((error) => {
        console.error('Error:', error);
    });
}).catch((error) => {
    console.error('Error:', error);
});

const time = new EventSource('/time');
const result = document.getElementById('result');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
let entryData = ""
workplace.textContent = sessionStorage.getItem("WorkplaceCode")
user.textContent = sessionStorage.getItem("User")
result.textContent = "Přihlaste se přiložením karty"

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

window.addEventListener("keyup", function (event) {
    entryData += event.key;
    if (entryData.includes("Enter")) {
        checkUser(entryData.toUpperCase());
        entryData = ""
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
                    OrderBarcode: sessionStorage.getItem("Order"),
                    DeviceId: sessionStorage.getItem("DeviceId"),
                    UserId: sessionStorage.getItem("UserId"),
                    Pcs: "0",
                    CloseLogin: "true"
                };
                fetch("/end_order", {
                    method: "POST",
                    body: JSON.stringify(data)
                }).then((response) => {
                    console.log("Ending order in Zapsi response: " + response.statusText);
                    sessionStorage.setItem("UserId", myObj.UserId)
                    sessionStorage.setItem("User", "Přihlášen: " + myObj.UserName)
                    data = {
                        OrderBarcode: sessionStorage.getItem("Order"),
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
                result.textContent = "Uživatel " + myObj.Result + " neexistuje v systému";
                setTimeout(() => result.textContent = "Přihlaste se přiložením karty", 3000)
            }
        });
    }).catch(() => {
        result.textContent = "Chyba komunikace";
        setTimeout(() => result.textContent = "Přihlaste se přiložením karty", 3000)
    });
}


window.focus()
