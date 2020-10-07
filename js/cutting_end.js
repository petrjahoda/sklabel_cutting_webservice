const time = new EventSource('/time');
const k2Pcs = document.getElementById('k2Pcs');
const enteredPcs = document.getElementById("enteredPcs");
const user = document.getElementById("user");
const btn0 = document.getElementById("0");
const btn1 = document.getElementById("1");
const btn2 = document.getElementById("2");
const btn3 = document.getElementById("3");
const btn4 = document.getElementById("4");
const btn5 = document.getElementById("5");
const btn6 = document.getElementById("6");
const btn7 = document.getElementById("7");
const btn8 = document.getElementById("8");
const btn9 = document.getElementById("9");
const btnDel = document.getElementById("del");
const btnOk = document.getElementById("ok");

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

let data = {Data: sessionStorage.getItem("Order")};
fetch("/get_k2Pcs", {
    method: "POST",
    body: JSON.stringify(data)
}).then((response) => {
    response.text().then(function (text) {
        let myObj = JSON.parse(text);
        k2Pcs.textContent = myObj.Data;
    });
}).catch((error) => {
    console.error('Error:', error);
});


function endOrder() {
    let data = {
        Type: "order",
        Code: "K302",
        WorkplaceCode: sessionStorage.getItem("WorkplaceCode"),
        UserId: sessionStorage.getItem("UserId"),
        OrderBarcode: sessionStorage.getItem("Order"),
        IdleId: sessionStorage.getItem("IdleId"),
        Pcs: enteredPcs.textContent
    };
    fetch("/save_code", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        console.log("Saving code to K2 response: " + response.statusText);
        let data = {
            OrderBarcode: sessionStorage.getItem("Order"),
            DeviceId: sessionStorage.getItem("DeviceId"),
            UserId: sessionStorage.getItem("UserId"),
            Pcs: enteredPcs.textContent
        };
        fetch("/end_order", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            console.log("Starting order in Zapsi response: " + response.statusText);
            window.location.replace('http://10.3.12:81/terminal/www/default/' + sessionStorage.getItem("DeviceId"));
        }).catch((error) => {
            console.error('Error:', error);
        });

    }).catch((error) => {
        console.error('Error:', error);
    });
}

btn0.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "0"
})
btn0.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "0"
})
btn1.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "1"
})
btn1.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "1"
})
btn2.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "2"
})
btn2.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "2"
})
btn3.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "3"
})
btn3.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "3"
})
btn4.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "4"
})
btn4.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "4"
})
btn5.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "5"
})
btn5.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "5"
})
btn6.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "6"
})
btn6.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "6"
})
btn7.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "7"
})
btn7.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "7"
})
btn8.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "8"
})
btn8.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "8"
})
btn9.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent + "9"
})
btn9.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent.slice(0, -1)
})
btnDel.addEventListener("click", () => {
    enteredPcs.textContent = enteredPcs.textContent.slice(0, -1)
})
btnDel.addEventListener("touchend", () => {
    enteredPcs.textContent = enteredPcs.textContent + "9"
})
btnOk.addEventListener("click", endOrder)
btnOk.addEventListener("touchend", endOrder)