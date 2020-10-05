const time = new EventSource('/time');
const result = document.getElementById('result');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

workplace.textContent = sessionStorage.getItem("WorkplaceCode")
user.textContent = sessionStorage.getItem("User")
result.textContent = "Probíhá řezání na zakázce číslo " + sessionStorage.getItem("Order")
