const time = new EventSource('/time');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
const result = document.getElementById('result');
const idleInput = document.getElementById('idle-input');
const cuttingEnd = document.getElementById('cutting-end');
const userChange = document.getElementById('user-change');
const userBreak = document.getElementById('user-break');

workplace.textContent = sessionStorage.getItem("WorkplaceCode")
user.textContent = sessionStorage.getItem("User")
result.textContent = "Probíhá řezání na zakázce číslo " + sessionStorage.getItem("Order")

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

idleInput.addEventListener("click", () => {
    window.location.replace('/idle_select')
})
idleInput.addEventListener("touchend", () => {
    window.location.replace('/idle_select')
})
cuttingEnd.addEventListener("click", () => {
    window.location.replace('/entry_pcs')
})
cuttingEnd.addEventListener("touchend", () => {
    window.location.replace('/entry_pcs')
})
userChange.addEventListener("click", () => {
    window.location.replace('/login')
})
userChange.addEventListener("touchend", () => {
    window.location.replace('/login')
})
userBreak.addEventListener("click", () => {
    window.location.replace('/login')
})
userBreak.addEventListener("touchend", () => {
    window.location.replace('/login')
})

