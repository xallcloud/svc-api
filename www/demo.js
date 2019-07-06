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

    document.getElementById('barra-cp-state').src="./guide-center-" + stateId + ".fw.png";
    
    // options were it can be 0 or 1
    var SecState = "0";
    if (stateId != "0") {
        SecState = "1"
    }
    document.getElementById('barra-cp-progress1').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-cp-progress2').src="./guide-line-" + SecState + ".fw.png";
}


function setStateApi(stateId) {
    console.log("setStateApi: ", stateId);

    document.getElementById('barra-api-state').src="./guide-center-" + stateId + ".fw.png";
    
    // options were it can be 0 or 1
    var SecState = "0";
    if (stateId != "0") {
        SecState = "1"
    }
    document.getElementById('barra-api-progress1').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-api-progress2').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-api-v').src="./guide-vline-" + SecState + ".fw.png";
}

function setStateDispatcher(stateId) {
    console.log("setStateDispatcher: ", stateId);

    document.getElementById('barra-dispatcher-state').src="./guide-center-" + stateId + ".fw.png";
    
    // options were it can be 0 or 1
    var SecState = "0";
    if (stateId != "0") {
        SecState = "1"
    }
    document.getElementById('barra-dispatcher-progress1').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-dispatcher-progress2').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-dispatcher-v').src="./guide-vline-" + SecState + ".fw.png";
}

function setStateNotify(stateId) {
    console.log("setStateNotify: ", stateId);

    document.getElementById('barra-notify-state').src="./guide-center-" + stateId + ".fw.png";
    
    // options were it can be 0 or 1
    var SecState = "0";
    if (stateId != "0") {
        SecState = "1"
    }
    document.getElementById('barra-notify-progress1').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-notify-progress2').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-notify-progress3').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-notify-progress4').src="./guide-line-" + SecState + ".fw.png";
    document.getElementById('barra-notify-v').src="./guide-vline-" + SecState + ".fw.png";
}

function setStateDevice1(stateId) {
    console.log("setStateDevice1: ", stateId);

    document.getElementById('barra-device1-state').src="./guide-center-" + stateId + ".fw.png";
}

function setStateDevice2(stateId) {
    console.log("setStateDevice2: ", stateId);

    document.getElementById('barra-device2-state').src="./guide-center-" + stateId + ".fw.png";
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
            
            setStateNotify("2");
        }


        if (activity.evDescription == "Attempting to reach end device.") {
            //console.log("GOT --> (Attempting to reach end device.)", activity.dvID);

            //setStateNotify("2");

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

        if (activity.evDescription == "User response: ack") {
            //console.log("GOT --> (User response: ack)", activity.dvID);

            if (activity.dvID == "UID-DEV-0000-0001") {
                setStateDevice1("2");
                FinalACK += 1;
            }

            if (activity.dvID == "UID-DEV-1000-0002") {
                setStateDevice2("2");
                FinalACK += 1;
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
        
        if (FinalACK > 0) {
            setStateCp("2");
        } else{
            setStateCp("3");
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
    var device_phone = document.getElementById('dev-phone')
    var device_pc = document.getElementById('dev-computer')
    var device_bell = document.getElementById('dev-bell')

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
                interval_controller = setInterval(poll, 2000)

                StartTime = "2019-07-05 22:56:23.000000 +0000 UTC";

                setStateApi("2");

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