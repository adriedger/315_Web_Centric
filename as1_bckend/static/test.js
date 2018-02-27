// CMPT 315 (Winter 2018)

//gets list of classes and runs showClasses
function getClasses() {
    let req = new XMLHttpRequest();

    req.addEventListener("load", function(evt) {
		if (req.response) {
			let data = JSON.parse(req.response);
			//console.log(data);
			showClasses(data);
		} else {
			console.log("no response");
		}
    });

    req.open("GET", `http://localhost:8080/api/v1/classes`);
    req.send();
}

function showClasses(data) {
    let container = document.querySelector("#class_list");

    // remove existing records
    while (container.childElementCount > 0) {
		container.removeChild(container.firstElementChild);
    }

    // add new records
    let table = document.createElement("table");

    if (data.length == 0) {
    	let row = document.createElement("tr");
		row.innerHTML = `No Classes in the System`
		table.appendChild(row);
    }
   
    for (let i = 0; i < data.length; i++) {
		let row = document.createElement("tr");
		row.innerHTML = `<td>${data[i].ClassName}</td>`
		table.appendChild(row);
    }

    container.appendChild(table);
}

//adds class, displays ok or error
function onCreateClick(evt) {
	//console.log(evt);
    let className = document.querySelector("#class_name_create").value;    
    //console.log(className);
    createClass(className);
}

function createClass(className) {
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			let data = JSON.parse(this.response);
			//console.log(data);
			//location.reload();
			getClasses();
			document.getElementById("create_response").innerHTML = `${data.ClassName} created. IMPORTANT. Creator key is: ${data.CreatorKey}`;
			document.getElementById("class_name_create").value = '';
		}
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/classes/create');
	//xhttp.setRequestHeader("Content-type", "application/json");
	//console.log(JSON.stringify({ClassName:`${className}`}));
	xhttp.send(JSON.stringify({ClassName:`${className}`}));
}

function onJoinClick(evt) {
    let className = document.querySelector("#class_name_join").value;
    let userName = document.querySelector("#user_name_join").value;
    joinClass(className, userName);
}

function joinClass(className, userName) {
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			let data = JSON.parse(this.response);
			getClasses();
			document.getElementById("join_response").innerHTML = `${data.Username} has joined class: ${data.ClassName}`;
			document.getElementById("class_name_join").value = '';
			document.getElementById("user_name_join").value = '';
		}
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/classes/join');
	xhttp.send(JSON.stringify({ClassName:`${className}`, Username:`${userName}`}));
}

function onQuestionsClick() {
	let className = document.querySelector("#class_name_questions").value;
	//console.log(className);
    getQuestions(className);
}

function getQuestions(className) {
    let req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
		if (req.response) {
			let data = JSON.parse(req.response);
			//console.log(data);
			showQuestions(data, className);
		} else {
			console.log("no response");
		}
    });
    req.open("GET", `http://localhost:8080/api/v1/classes/questions/${className}`);
    req.send();
}

function showQuestions(data, className) {
    let container = document.querySelector("#question_list");
    while (container.childElementCount > 0) {
		container.removeChild(container.firstElementChild);
    }

    var head = document.createElement("h");
    var t = document.createTextNode(`Class ${className}`);
    head.appendChild(t);
    container.appendChild(head)

    let table = document.createElement("table");
    if (data.length == 0) {
    	let row = document.createElement("tr");
		row.innerHTML = `No Questions in this Class`
		table.appendChild(row);
    }
    for (let i = 0; i < data.length; i++) {
		let row = document.createElement("tr");
		row.innerHTML = `<td>${data[i].Question}</td>`
		table.appendChild(row);
    }
    container.appendChild(table);
}

//checks for button presses
function addEventListeners() {
    let createButton = document.querySelector("#create_button");
    createButton.addEventListener("click", onCreateClick);

    let joinButton = document.querySelector("#join_button");
    joinButton.addEventListener("click", onJoinClick);

    let questionsButton = document.querySelector("#questions_button");
    questionsButton.addEventListener("click", onQuestionsClick);
}

getClasses();
addEventListeners();