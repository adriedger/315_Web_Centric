// Andre Driedger

function getClasses() {
    let req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
		if (req.response) {
			let data = JSON.parse(req.response);
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
    container.innerHTML = "Classes"
    while (container.childElementCount > 0) {
		container.removeChild(container.firstElementChild);
    }
    let table = document.createElement("table");
    if (data.length == 0) {
    	container.innerHTML = "No Classes in the System";
    }   
    for (let i = 0; i < data.length; i++) {
		let row = document.createElement("tr");
		row.innerHTML = `<td>${data[i].ClassName}</td>`
		table.appendChild(row);
    }
    container.appendChild(table);
}

function onCreateClick(evt) {
    let className = document.querySelector("#class_name_create").value;    
    createClass(className);
}

function createClass(className) {
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			let data = JSON.parse(this.response);
			getClasses();
			document.getElementById("create_response").innerHTML = `${data.ClassName} created. IMPORTANT. Creator key is: ${data.CreatorKey}`;			
		} else {
			document.getElementById("create_response").innerHTML = `UNSUCCESSFUL: ${className} already exists`;
		}
		document.getElementById("class_name_create").value = '';
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/classes/create');
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
			document.getElementById("join_response").innerHTML = `${data.Username} has joined ${data.ClassName}`;
		} else {
			document.getElementById("join_response").innerHTML = `UNSUCCESSFUL: ${userName} is already in ${className}`;
		}
		document.getElementById("class_name_join").value = '';
		document.getElementById("user_name_join").value = '';
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/classes/join');
	xhttp.send(JSON.stringify({ClassName:`${className}`, Username:`${userName}`}));
}

function onQuestionsClick() {
	document.getElementById("general_responses").innerHTML = "";
	document.getElementById("class_key").value = "";	
	document.getElementById("question_viewresponses").value = "";
	document.getElementById("username_student").value = "";
	document.getElementById("responses_list").innerHTML = "";

	let o = document.querySelector("#option_selector");
	o.style.display = "block";
	
	let className = document.querySelector("#class_name_questions").value;
    getQuestions(className, false);
}

function getQuestions(className, creator) {
    let req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
		if (req.response) {
			let data = JSON.parse(req.response);
			showQuestions(data, creator);
		} else {
			console.log("no response");
		}
    });
    req.open("GET", `http://localhost:8080/api/v1/classes/questions/${className}`);
    req.send();
}

function showQuestions(data, creator) {
    let container = document.querySelector("#question_list");
    container.innerHTML = "Questions"
    while (container.childElementCount > 0) {
		container.removeChild(container.firstElementChild);
    }
    let table = document.createElement("table");
    if (data.length == 0) {
    	container.innerHTML = "No Questions in this Class"
    }
    for (let i = 0; i < data.length; i++) {
    	let row = document.createElement("tr");
    	if (creator) {			
			row.innerHTML = `<td>Question: ${data[i].Question} Answer: ${data[i].Answer}</td>`
    	} else {
			row.innerHTML = `<td>${data[i].Question}</td>`
    	}
		table.appendChild(row);
    }
    container.appendChild(table);
}

function onStudentClick() {
	let className = document.querySelector("#class_name_questions").value;
	getQuestions(className, false);
	document.getElementById("general_responses").innerHTML = "";
	let s = document.querySelector("#student_options");
	let c = document.querySelector("#creator_options");
	let o = document.querySelector("#creator_stuff");
	c.style.display = "none";
	o.style.display = "none";
	s.style.display = "block";
}

function onCreatorClick() {
	document.getElementById("general_responses").innerHTML = "";
	let s = document.querySelector("#student_options");
	let c = document.querySelector("#creator_options");
	s.style.display = "none";
	c.style.display = "block";
}

function onVerifyCreatorClick() {
	let className = document.querySelector("#class_name_questions").value;
	let classKey = document.querySelector("#class_key").value;
	verifyCreator(className, classKey)
}

function verifyCreator(className, classKey) {
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
  		let c = document.querySelector("#creator_stuff");
		if (this.readyState == 4 && this.status == 200) {
			getQuestions(className, true);
			c.style.display = "block";
			document.getElementById("general_responses").innerHTML = `Key verified`;			
		} else {
			c.style.display = "none";
			document.getElementById("general_responses").innerHTML = `UNSUCCESSFUL: Key does not match class`;
		}
		document.getElementById("question_addquestion").value = '';
		document.getElementById("answer_addquestion").value = '';
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/questions/responses');
	xhttp.send(JSON.stringify({ClassName:`${className}`,KeyAttempt:`${classKey}`,Question:``}));
}

function onAddQuestionClick() {
	let className = document.querySelector("#class_name_questions").value;
	let classKey = document.querySelector("#class_key").value;
	let question = document.querySelector("#question_addquestion").value;
	let answer = document.querySelector("#answer_addquestion").value;
	addQuestion(className, classKey, question, answer);
}

function addQuestion(className, classKey, question, answer) {
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			getQuestions(className, true);
			document.getElementById("general_responses").innerHTML = `Question ${question} created. Answer is ${answer}`;			
		} else {
			document.getElementById("general_responses").innerHTML = `UNSUCCESSFUL: question not created`;
		}
		document.getElementById("question_addquestion").value = '';
		document.getElementById("answer_addquestion").value = '';
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/questions/create');
	xhttp.send(JSON.stringify({ClassName:`${className}`,KeyAttempt:`${classKey}`,Question:`${question}`,Answer:`${answer}`}));	
}

function onDeleteQuestionClick() {
	let className = document.querySelector("#class_name_questions").value;
	let classKey = document.querySelector("#class_key").value;
	let question = document.querySelector("#question_deletequestion").value;
	deleteQuestion(className, classKey, question);
}

function deleteQuestion(className, classKey, question) {
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			getQuestions(className, true);
			document.getElementById("general_responses").innerHTML = `Question ${question} deleted`;			
		} else {
			document.getElementById("general_responses").innerHTML = `UNSUCCESSFUL: Question ${question} not deleted`;
		}
		document.getElementById("question_deletequestion").value = '';
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/questions/delete');
	xhttp.send(JSON.stringify({ClassName:`${className}`,KeyAttempt:`${classKey}`,Question:`${question}`}));
}

function onViewStudentResponsesClick() {
	let className = document.querySelector("#class_name_questions").value;
	let classKey = document.querySelector("#class_key").value;
	let question = document.querySelector("#question_viewresponses").value;
	getStudentResponses(className, classKey, question);
}

function getStudentResponses(className, classKey, question) {
	let req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
		if (req.response) {
			let data = JSON.parse(req.response);
			showStudentResponses(data);
		} else {
			console.log("no response");
		}
    });
    req.open("POST", `http://localhost:8080/api/v1/questions/responses`);
    req.send(JSON.stringify({ClassName:`${className}`,KeyAttempt:`${classKey}`,Question:`${question}`}));
}

function showStudentResponses(data) {
	let container = document.querySelector("#responses_list");
	container.innerHTML = "Responses"
    while (container.childElementCount > 0) {
		container.removeChild(container.firstElementChild);
    }
    let table = document.createElement("table");
    if (data.length == 0) {
    	container.innerHTML = "No Student Responses for this Question"
    }
    for (let i = 0; i < data.length; i++) {
    	let row = document.createElement("tr");		
		row.innerHTML = `<td>${data[i].Username}: ${data[i].Response}</td>`
		table.appendChild(row);
    }
    container.appendChild(table);
    console.log(data)
}

function onAddResponseClick() {
	let className = document.querySelector("#class_name_questions").value;
	let userName = document.querySelector("#username_student").value;
	let question = document.querySelector("#question_addresponse").value;
	let response = document.querySelector("#response_addresponse").value;
	addResponse(className, userName, question, response);
}

function addResponse(className, userName, question, response) {
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			document.getElementById("general_responses").innerHTML = `Answer: ${response} Submitted for Question: ${question}`;			
		} else {
			document.getElementById("general_responses").innerHTML = `UNSUCCESSFUL: Response not submitted`;
		}
		document.getElementById("question_addresponse").value = '';
		document.getElementById("response_addresponse").value = '';
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/responses/add');
	xhttp.send(JSON.stringify({ClassName:`${className}`,Username:`${userName}`,Question:`${question}`,Response:`${response}`}));
}

function onModifyResponseClick() {
	let className = document.querySelector("#class_name_questions").value;
	let userName = document.querySelector("#username_student").value;
	let question = document.querySelector("#question_modifyresponse").value;
	let response = document.querySelector("#response_modifyresponse").value;
	modifyResponse(className, userName, question, response);
}

function modifyResponse(className, userName, question, response){
	var xhttp = new XMLHttpRequest();
  	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			document.getElementById("general_responses").innerHTML = `Answer: ${response} Updated for Question: ${question}`;			
		} else {
			document.getElementById("general_responses").innerHTML = `UNSUCCESSFUL: Response not modified`;
		}
		document.getElementById("question_modifyresponse").value = '';
		document.getElementById("response_modifyresponse").value = '';
	};
	xhttp.open('POST', 'http://localhost:8080/api/v1/responses/modify');
	xhttp.send(JSON.stringify({ClassName:`${className}`,Username:`${userName}`,Question:`${question}`,Response:`${response}`}));
}

function addEventListeners() {
    let createButton = document.querySelector("#create_button");
    createButton.addEventListener("click", onCreateClick);

    let joinButton = document.querySelector("#join_button");
    joinButton.addEventListener("click", onJoinClick);

    let questionsButton = document.querySelector("#questions_button");
    questionsButton.addEventListener("click", onQuestionsClick);

    let studentButton = document.querySelector("#student_button");
    studentButton.addEventListener("click", onStudentClick);

    let creatorButton = document.querySelector("#creator_button");
    creatorButton.addEventListener("click", onCreatorClick);

    let addQuestionButton = document.querySelector("#addquestion_button");
    addQuestionButton.addEventListener("click", onAddQuestionClick);

    let deleteQuestionButton = document.querySelector("#deletequestion_button");
    deleteQuestionButton.addEventListener("click", onDeleteQuestionClick);

    let viewResponsesButton = document.querySelector("#viewresponses_button");
    viewResponsesButton.addEventListener("click", onViewStudentResponsesClick);

    let addResponseButton = document.querySelector("#addresponse_button");
    addResponseButton.addEventListener("click", onAddResponseClick);

    let modifyResponseButton = document.querySelector("#modifyresponse_button");
    modifyResponseButton.addEventListener("click", onModifyResponseClick);

    let creatorVerifyButton = document.querySelector("#viewanswers_button");
    creatorVerifyButton.addEventListener("click", onVerifyCreatorClick);
}

getClasses();
addEventListeners();