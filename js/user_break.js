let data = {
    Order: sessionStorage.getItem("Order"),
    DeviceId: sessionStorage.getItem("DeviceId"),
    UserId: sessionStorage.getItem("UserId")
};
fetch("/end_order", {
    method: "POST",
    body: JSON.stringify(data)
}).then((response) => {
    console.log("Ending order in Zapsi response: " + response.statusText);
    let data = {Data: "K219"};
    fetch("/save_code", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        console.log("Saving code to K2 response: " + response.statusText);
        let data = {Data: "0004"};
        fetch("/save_code", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            console.log("Saving code to K2 response: " + response.statusText);
            sessionStorage.setItem("UserId", "")
            sessionStorage.setItem("User", "Přihlášen: " + "")
            user.textContent = "Nenávaznost obsluhy"
            let data = {
                Order: sessionStorage.getItem("Order"),
                DeviceId: sessionStorage.getItem("DeviceId"),
                UserId: sessionStorage.getItem("")
            };
            fetch("/start_order", {
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
    entryData = entryData.replace("Meta", "");
    entryData = entryData.replaceAll("Enter", "");
    entryData = entryData.replaceAll("Shift", "");
    if (event.key === "Enter" && entryData.length > 0) {
        console.log("DATA: " + entryData)
        checkUser(entryData);
        entryData = "";
    } else {
        entryData += event.key;
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
                let result = ""
                fetch("/end_order", {
                    method: "POST",
                    body: JSON.stringify(data)
                }).then((response) => {
                    console.log("Ending order in Zapsi response: " + response.statusText);
                    sessionStorage.setItem("UserId", myObj.UserId)
                    sessionStorage.setItem("User", "Přihlášen: " + myObj.UserName)
                    let data = {
                        Order: barcode,
                        DeviceId: sessionStorage.getItem("DeviceId"),
                        UserId: sessionStorage.getItem("UserId")
                    };
                    fetch("/start_order", {
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
                result.textContent = "Uživatel " + barcode + " neexistuje v systému";
                setTimeout(() => result.textContent = "Přihlaste se přiložením karty", 3000)
            }
        });
    }).catch((error) => {
        console.error('Error:', error);
        result.textContent = "Problém v komunikaci se serverem.";
    });
}



