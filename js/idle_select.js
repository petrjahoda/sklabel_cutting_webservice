const time = new EventSource('/time');
const result = document.getElementById('result');
const workplace = document.getElementById("workplace");
const user = document.getElementById("user");
const previous = document.getElementById("previous");
const next = document.getElementById("next");
let idles
let actualPage = 1
let totalPages = 1

workplace.textContent = sessionStorage.getItem("WorkplaceCode")
user.textContent = sessionStorage.getItem("User")

time.addEventListener('time', (e) => {
    document.getElementById("time").innerHTML = e.data;
}, false);

fetch("/get_idles", {
    method: "POST",

}).then((response) => {
    response.text().then(function (text) {
        let myObj = JSON.parse(text);
        idles = myObj.Data;
        totalPages = Math.ceil(idles.length / 12)
        let layer = 4;
        let position = 1;
        let numberOfIdles = 0;
        idles.forEach(function (idle) {
            let btn = document.createElement("BUTTON");
            btn.innerHTML = idle.Name;
            btn.dataset.id = idle.Id;
            btn.classList.add("button");
            btn.classList.add("cnt3");
            btn.classList.add("pos" + position);
            btn.classList.add("layer" + layer);
            btn.classList.add("yellow");
            btn.id = numberOfIdles + 1;
            position++
            numberOfIdles++
            if (position > 3) {
                position = 1
                layer--
            }
            if (numberOfIdles <= 12) {
                document.body.appendChild(btn);
                btn.addEventListener("click", function (event) {
                    console.log(btn.dataset.id)
                    //TODO: send K2
                    //TODO: create Idle
                    //TODO: show "idle is running" screen
                });
                btn.addEventListener("touchend", function(event) {

                });
            } else {
                next.style.display = "block";
                actualPage = 1
            }
        });

    });
}).catch((error) => {
    console.error('Error:', error);
});


function showPrevious() {
    actualPage = actualPage - 1
    for (let i = 4; i > 0; i--) {
        let actualIdles = document.querySelectorAll(".layer" + i)
        console.log(actualIdles.length)
        actualIdles.forEach(function (idle) {
            idle.remove()
        });
    }
    let layer = 4;
    let position = 1;
    let numberOfIdles = 0
    idles.forEach(function (idle) {
        let btn = document.createElement("BUTTON");
        btn.innerHTML = idle.Name;
        btn.dataset.id = idle.Id;
        btn.classList.add("button");
        btn.classList.add("cnt3");
        btn.classList.add("pos" + position);
        btn.classList.add("layer" + layer);
        btn.classList.add("yellow");
        btn.id = numberOfIdles + 1;
        position++
        numberOfIdles++
        if (position > 3) {
            position = 1
            layer--
        }
        if (layer === 0) {
            layer = 4
        }
        if (numberOfIdles > 12 * (actualPage - 1) && numberOfIdles <= (12 * (actualPage))) {
            document.body.appendChild(btn);
            // btn.addEventListener("click", idleSelected(this));
            // btn.addEventListener("touch", idleSelected(this));
        }
        if (numberOfIdles < (12 * actualPage)) {
            next.style.display = "none"
        } else {
            next.style.display = "block"
        }
        if (actualPage === 1) {
            previous.style.display = "none"
        }
    });
}

// function idleSelected() {
//     console.log()
// }

function showNext() {
    previous.style.display = "block"
    actualPage = actualPage + 1
    for (let i = 4; i > 0; i--) {
        let actualIdles = document.querySelectorAll(".layer" + i)
        console.log(actualIdles.length)
        actualIdles.forEach(function (idle) {
            idle.remove()
        });
    }
    let layer = 4;
    let position = 1;
    let numberOfIdles = 0
    idles.forEach(function (idle) {
        let btn = document.createElement("BUTTON");
        btn.innerHTML = idle.Name;
        btn.dataset.id = idle.Id;
        btn.classList.add("button");
        btn.classList.add("cnt3");
        btn.classList.add("pos" + position);
        btn.classList.add("layer" + layer);
        btn.classList.add("yellow");
        btn.id = numberOfIdles + 1;
        position++
        numberOfIdles++
        if (position > 3) {
            position = 1
            layer--
        }
        if (layer === 0) {
            layer = 4
        }
        if (numberOfIdles > 12 * (actualPage - 1) && numberOfIdles <= (12 * (actualPage))) {
            document.body.appendChild(btn);
            // btn.addEventListener("click", idleSelected(this));
            // btn.addEventListener("touch", idleSelected(this));
        }
        if (numberOfIdles > (12 * actualPage)) {
            next.style.display = "block"
        } else {
            next.style.display = "none"
        }
    });
}

previous.addEventListener("click", showPrevious)
previous.addEventListener("touchend", showPrevious)
next.addEventListener("click", showNext)
next.addEventListener("touchend", showNext)