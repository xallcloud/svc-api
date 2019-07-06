var ACID = ""
var StartTime = ""
var interval_controller = 0
var activity_list = []
var device_list = [
    {
        id: "UID-DEV-0000-0001",
        name: "Phone 1"
    },
    {
        id: "UID-DEV-1000-0002",
        name: "Mobile 2"
    }
]

/**
 * Activa e desactiva o CPT.
 */
function toggleCPT() {
    var cpt = document.getElementById('cpt');

    if (cpt.classList.contains('activated')) {
        cpt.classList.remove('activated');
        clearInterval(interval_controller)
        ACID = ""
        
        setStateCp("0");
        setStateApi("0");
        setStateDispatcher("0");
        setStateNotify("0");
        setStateDevice1("0");
        setStateDevice2("0");
    } else {
        cpt.classList.add('activated');

        setStateCp("1");

        postCPT();
    }
}


function setStateCp(stateId) {
    console.log("setStateCp: ", stateId);

    if (stateId == "0") {
        document.getElementById('cpt').src = "./callpoint-0.png";
        document.getElementById('cp-in').src="./dev-out-0.png";
    } else if (stateId == "1") {
        document.getElementById('cp-in').src="./dev-out-1.png";
    } else if (stateId == "2") {
        document.getElementById('cp-in').src="./dev-out-2.png";
    } else if (stateId == "3") {
        document.getElementById('cpt').src = "./callpoint-3.png";
    } else if (stateId == "4") {
        document.getElementById('cpt').src = "./callpoint-4.png";
    } else if (stateId == "5") {
        document.getElementById('cpt').src = "./callpoint-5.png";
    }
}


function setStateApi(stateId) {
    console.log("setStateApi: ", stateId);

    if (stateId == "0") {
        document.getElementById('svc-api-top-out').src = "./svc-top-p2.png";
        document.getElementById('svc-api-bottom-out').src = "./svc-bottom-p0.png";
    } else {
        document.getElementById('svc-api-top-out').src = "./svc-top-out.png";
        document.getElementById('svc-api-bottom-out').src = "./svc-bottom-out.png";
    }
}

function setStateDispatcher(stateId) {
    console.log("setStateDispatcher: ", stateId);

    if (stateId == "0") {
        document.getElementById('svc-dispatcher-top-in').src = "./svc-top-p1.png";
        document.getElementById('svc-dispatcher-top-out').src = "./svc-top-p2.png";
        document.getElementById('svc-dispatcher-bottom-out').src = "./svc-bottom-p0.png";
    } else if (stateId == "1") {
        document.getElementById('svc-dispatcher-top-in').src = "./svc-top-in.png";
        document.getElementById('svc-dispatcher-top-out').src = "./svc-top-p2.png";
        document.getElementById('svc-dispatcher-bottom-out').src = "./svc-bottom-out.png";
    } else {
        document.getElementById('svc-dispatcher-top-in').src = "./svc-top-in.png";
        document.getElementById('svc-dispatcher-top-out').src = "./svc-top-out.png";
        document.getElementById('svc-dispatcher-bottom-out').src = "./svc-bottom-out.png";
    }
}

function setStateNotify(stateId) {
    console.log("setStateNotify: ", stateId);

    if (stateId == "0") {
        document.getElementById('svc-notify-top-in').src = "./svc-top-p1.png";
        document.getElementById('svc-notify-bottom-out').src = "./svc-bottom-p0.png";
    } else if (stateId == "1") {
        document.getElementById('svc-notify-top-in').src = "./svc-top-in.png";
        document.getElementById('svc-notify-bottom-out').src = "./svc-bottom-out.png";
    } else {
        document.getElementById('svc-notify-top-in').src = "./svc-top-in.png";
        document.getElementById('svc-notify-bottom-out').src = "./svc-bottom-out.png";
    }
}

function setStateDevice1(stateId) {
    console.log("setStateDevice1: ", stateId);

    document.getElementById('dev1-out').src="./dev-out-" + stateId + ".fw.png";

    if (stateId == "0") { //reset
        document.getElementById('dev1-out').src = "./dev-out-0.png";
        document.getElementById('dev1-in').src = "./dev1-in-0.png";
    } else if (stateId == "1") { //sending
        document.getElementById('dev1-out').src = "./dev-out-1.png";
        document.getElementById('dev1-in').src = "./dev1-in-0.png";
    } else if (stateId == "2") { //sent
        document.getElementById('dev1-out').src = "./dev-out-2.png";
        document.getElementById('dev1-in').src = "./dev1-in-0.png";
    } else if (stateId == "3") { //error sending
        document.getElementById('dev1-out').src = "./dev-out-3.png";
        document.getElementById('dev1-in').src = "./dev1-in-0.png";
    } else if (stateId == "4") { //acepted
        document.getElementById('dev1-out').src = "./dev-out-2.png";
        document.getElementById('dev1-in').src = "./dev1-in-4.png";
    } else { //user rejected
        document.getElementById('dev1-out').src = "./dev-out-2.png";
        document.getElementById('dev1-in').src = "./dev1-in-3.png";
    }
}

function setStateDevice2(stateId) {
    console.log("setStateDevice1: ", stateId);

    document.getElementById('dev2-out').src="./dev-out-" + stateId + ".fw.png";

    if (stateId == "0") { //reset
        document.getElementById('dev2-out').src = "./dev-out-0.png";
        document.getElementById('dev2-in').src = "./dev2-in-0.png";
    } else if (stateId == "1") { //sending
        document.getElementById('dev2-out').src = "./dev-out-1.png";
        document.getElementById('dev2-in').src = "./dev2-in-0.png";
    } else if (stateId == "2") { //sent
        document.getElementById('dev2-out').src = "./dev-out-2.png";
        document.getElementById('dev2-in').src = "./dev2-in-0.png";
    } else if (stateId == "3") { //error sending
        document.getElementById('dev2-out').src = "./dev-out-3.png";
        document.getElementById('dev2-in').src = "./dev2-in-0.png";
    } else if (stateId == "4") { //acepted
        document.getElementById('dev2-out').src = "./dev-out-2.png";
        document.getElementById('dev2-in').src = "./dev2-in-4.png";
    } else { //user rejected
        document.getElementById('dev2-out').src = "./dev-out-2.png";
        document.getElementById('dev2-in').src = "./dev2-in-3.png";
    }
}

/** Serve para gerar UID unicos */
function uuidv4() {
    return 'UID-xxxx-xxxx-xxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

/**
 * update HTML table
 */

function AddFirstRow(){

    var table = document.getElementById('table-body')
    table.innerHTML = "";

    var row = document.createElement('tr')
    var cell_activity = document.createElement('td')
    var cell_timestamp = document.createElement('td')
    var cell_device = document.createElement('td')

    cell_activity.innerText = "Action sent to API."
    cell_timestamp.innerText = StartTime
    cell_device.innerText = "--"

    row.appendChild(cell_device)
    row.appendChild(cell_activity)
    row.appendChild(cell_timestamp)

    table.appendChild(row)
}


function drawTable() {
    var table = document.getElementById('table-body')
    table.innerHTML = "";
    AddFirstRow();


    var FinalACK = 0;
    var FinalReject = 0;
    var FinalFailed = 0;
    var FinalState = 0;

    activity_list.forEach(function (activity) {

        console.log("activity.evDescription: ", activity.evDescription);

        // draw table
        if (activity.evDescription == "Notification started") {
            //console.log("GOT --> (Notification started)", activity.evDescription);
            
            //setStateCp("2");
            setStateApi("2");
            setStateDispatcher("1");
        }

        if (activity.evDescription == "Notification sent to svc-notify.") {
            //console.log("GOT --> (Notification sent to svc-notify.)", activity.evDescription);
            
            setStateDispatcher("2");
            setStateNotify("1");
        }


        if (activity.evDescription == "Got Notification on svc-notify.") {
            //console.log("GOT --> (Notification sent to svc-notify.)", activity.evDescription);
            
            setStateNotify("1");
        }


        if (activity.evDescription == "Attempting to reach end device.") {
            //console.log("GOT --> (Attempting to reach end device.)", activity.dvID);

            setStateNotify("2");

            if (activity.dvID == "UID-DEV-0000-0001") {
                setStateDevice1("1");
            }

            if (activity.dvID == "UID-DEV-1000-0002") {
                setStateDevice2("1");
            }        
        }

        if (activity.evDescription == "Failed to deliver message to end device.") {
            //console.log("GOT --> (Failed to deliver message to end device.)", activity.dvID);

            if (activity.dvID == "UID-DEV-0000-0001") {
                setStateDevice1("3");
                FinalFailed += 1;
            }

            if (activity.dvID == "UID-DEV-1000-0002") {
                setStateDevice2("3");
                FinalFailed += 1;
            }        
        }

        if (activity.evDescription == "Timeout trying to deliver message to end device.") {
            //console.log("GOT --> (Failed to deliver message to end device.)", activity.dvID);

            if (activity.dvID == "UID-DEV-0000-0001") {
                setStateDevice1("3");
                FinalFailed += 1;
            }

            if (activity.dvID == "UID-DEV-1000-0002") {
                setStateDevice2("3");
                FinalFailed += 1;
            }        
        }

        if (activity.evDescription == "User response: ack") {
            //console.log("GOT --> (User response: ack)", activity.dvID);

            if (activity.dvID == "UID-DEV-0000-0001") {
                setStateDevice1("4");
                setStateCp("4");
                FinalACK += 1;
            }

            if (activity.dvID == "UID-DEV-1000-0002") {
                setStateDevice2("4");
                setStateCp("4");
                FinalACK += 1;
            }        
        }

        if (activity.evDescription == "User response: cancel") {
            //console.log("GOT --> (User response: ack)", activity.dvID);

            if (activity.dvID == "UID-DEV-0000-0001") {
                setStateDevice1("5");
                //setStateCp("5");
                FinalReject += 1;
            }

            if (activity.dvID == "UID-DEV-1000-0002") {
                setStateDevice2("5");
                //setStateCp("5");
                FinalReject += 1;
            }        
        }

        if (activity.evDescription == "Notification reached final state.") {
            //console.log("GOT --> (User response: ack)", activity.dvID);            
            FinalState += 1;            
        }

        

        
        var row = document.createElement('tr')

        var cell_activity = document.createElement('td')
        var cell_timestamp = document.createElement('td')
        var cell_device = document.createElement('td')

        cell_activity.innerText = activity.evDescription
        cell_timestamp.innerText = activity.created
        cell_device.innerText = getDeviceName(activity.dvID)

        row.appendChild(cell_device)
        row.appendChild(cell_activity)
        row.appendChild(cell_timestamp)

        table.appendChild(row)
        
    })


    console.log("GOT --> FinalACK", FinalACK);
    console.log("GOT --> FinalFailed", FinalFailed);

    if (FinalState >= 2){
        console.log("Final State", FinalACK);
        
        if (FinalACK == 0) {
            if (FinalReject >= 1){
                setStateCp("5");            
            } else {
                setStateCp("3"); 
            }

        }
        //Stop
        clearInterval(interval_controller)
    }
}

/**
 *  Get Device Name by UID
 * @param {string} deviceId 
 */
function getDeviceName(deviceId) {
    console.log(deviceId, device_list)
    var found = device_list.find(function (device) {
        return device.id == deviceId
    })

    if (found) {
        return found.name
    } else {
        return "";
    }
}

function poll() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "https://xall.cloud/api/events/action/" + ACID, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    //xhr.setRequestHeader('Access-Control-Allow-Origin', '*');

    xhr.onreadystatechange = function () {
        if (this.readyState != 4) return;

        if (this.status < 300) {
            var data = JSON.parse(this.responseText);
            activity_list = data

            processActivities()

        }
    };

    xhr.send();
}

/**
 * Process activities
 */
function processActivities() {
    /** Do stuff with "activity_list" */


    /** devices */
    //var device_phone = document.getElementById('dev-phone')
    //var device_pc = document.getElementById('dev-computer')
    //var device_bell = document.getElementById('dev-bell')

    // add CSS class:
    //  device_phone.classList.add('RED')   --> classes implementadas "RED", "GREEN", "BLUE", "GREY"

    drawTable();
}


function postCPT() {

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "https://xall.cloud/api/action", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    //xhr.setRequestHeader('Access-Control-Allow-Origin', '*');

    xhr.onreadystatechange = function () {
        if (this.readyState != 4) return;

        if (this.status < 300) {
            var data = JSON.parse(this.responseText);
            if (data.AcID && data.KeyID) {
                interval_controller = setInterval(poll, 1000)

                StartTime = "2019-07-05 22:56:23.000000 +0000 UTC";

                setStateApi("1");
                setStateCp("2");

                AddFirstRow();

            }
            console.log("GOT --> ", data)
        }
    };

    ACID = uuidv4()

    xhr.send(JSON.stringify({
        "acID": ACID,
        "cpID": "UID-1000-0000-0001",
        "action": "activate",
        "description": "Activate Callpoint"
    }));

    

}