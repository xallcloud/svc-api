var ACID = ""
var interval_controller = 0
var activity_list = []
var device_list = [
    {
        id: "UID-DEV-0000-0001",
        name: "Telefone do Mário"
    },
    {
        id: "UID-DEV-1000-0002",
        name: "Computador do Isac"
    },
    {
        id: "UID-DEV-1000-0003",
        name: "Campainha do Tito"
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
    } else {
        cpt.classList.add('activated');
        postCPT();
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
function drawTable() {
    var table = document.getElementById('table-body')
    table.innerHTML = "";
    activity_list.forEach(function (activity) {
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
    xhr.setRequestHeader('Access-Control-Allow-Origin', '*');

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
    xhr.setRequestHeader('Access-Control-Allow-Origin', '*');

    xhr.onreadystatechange = function () {
        if (this.readyState != 4) return;

        if (this.status < 300) {
            var data = JSON.parse(this.responseText);
            if (data.AcID && data.KeyID) {
                interval_controller = setInterval(poll, 2000)
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