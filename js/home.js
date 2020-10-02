const time = new EventSource('/time');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
const result = document.getElementById('result');

workplace.textContent = sessionStorage.getItem("WorkplaceCode")
user.textContent = sessionStorage.getItem("User")
result.textContent = "Probíhá řezání na zakázce číslo " + sessionStorage.getItem("Order")

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);