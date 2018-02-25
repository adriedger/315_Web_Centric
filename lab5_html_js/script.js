// CMPT 315 (Winter 2018)
// Lab #5: Sample Solution

function showSchedule(sched) {
    let container = document.querySelector("#schedule");

    // remove existing records
    while (container.childElementCount > 0) {
	container.removeChild(container.firstElementChild);
    }

    // add new records
    let records = sched.payload;
    let table = document.createElement("table");
    for (let i = 0; i < records.length; i++) {
	let row = document.createElement("tr");
	row.innerHTML = `
<td>${records[i].r}</td>
<td>${records[i].s}</td>
<td>${records[i].e}</td>`
	table.appendChild(row);
    }

    container.appendChild(table);
}

function getSchedule(stop) {
    let req = new XMLHttpRequest();

    req.addEventListener("load", function(evt) {
	if (req.response) {
	    let data = JSON.parse(req.response);
	    showSchedule(data);
	} else {
	    console.log("no response");
	}
    });

    req.open("GET", `https://smartbus.ca/stop/${stop}`);
    req.send();
}

function onGetScheduleClick(evt) {
    let stop = document.querySelector("#stopNumber").value;
    
    getSchedule(stop);
}

function addEventListeners() {
    let btn = document.querySelector("#getSchedule");

    btn.addEventListener("click", onGetScheduleClick);
}

addEventListeners();
