// CMPT 315 (Winter 2018)

function showClasses(data) {
    let container = document.querySelector("#class_list");

    // remove existing records
    while (container.childElementCount > 0) {
	container.removeChild(container.firstElementChild);
    }

    // add new records
    let table = document.createElement("table");
   
    for (let i = 0; i < data.length; i++) {
		let row = document.createElement("tr");
		row.innerHTML = `<td>${data[i].ClassName}</td>`
		table.appendChild(row);
    }

    container.appendChild(table);
}

//gets list of classes and runs showClasses
function getClasses() {
    let req = new XMLHttpRequest();

    req.addEventListener("load", function(evt) {
	if (req.response) {
	    let data = JSON.parse(req.response);
	    console.log(data);
	    showClasses(data);
	} else {
	    console.log("no response");
	}
    });

    req.open("GET", `http://localhost:8080/api/v1/classes`);
    req.send();
}

//adds class, displays ok or error
function onGetCreateClick(evt) {
	//console.log(evt);
    let className = document.querySelector("#class_name_create").value;    
    //console.log(className);
    createClass(className);
}

function createClass(className) {
	var xhttp = new XMLHttpRequest();
//  	xhttp.onreadystatechange = function() {
//    if (this.readyState == 4 && this.status == 200) {
//		document.getElementById("demo").innerHTML = this.responseText;
//		}
//	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/classes/create');
	xhttp.setRequestHeader("Content-type", "application/json");
	console.log(JSON.stringify({ClassName:`${className}`}));
	xhttp.send(JSON.stringify({ClassName:`${className}`}));
}

//checks for button presses
function addEventListeners() {
    let createButton = document.querySelector("#create_button");
    createButton.addEventListener("click", onGetCreateClick);
}

getClasses();
addEventListeners();