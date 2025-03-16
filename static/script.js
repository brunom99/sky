let map;
let socket;
let tryToConnecting = false;
let aircrafts = {};
const lengthMaxLines = 100;

function gID(id) {
    return document.getElementById(id);
}

function showError(err) {
    gID("error").innerHTML = err;
}

async function fetchAPI(apiName) {
    try {
        const response = await fetch("api/" + apiName);
        if (!response.ok) {
            throw new Error("response is not ok");
        }
        return await response.json();
    } catch (error) {
        showError(apiName + ": " + error);
    }
}

function onload() {
    initMap();
    lastActivity();
}

function initMap() {
    // init map
    map = L.map('map');
    // add copyright
    L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);
    // target Geneve
    var geneve = L.latLng('46.204391', '6.143158');
    map.setView(geneve, 2);
}

async function lastActivity() {
    const activity = await fetchAPI("activity");
    const divActivity = gID("activity");
    if (activity) {
        divActivity.innerHTML = 'last server activity: ' + new Date(activity.last_activity).toTimeString().slice(0, 8);
    } else {
        divActivity.innerHTML = '';
    }
    setTimeout(lastActivity, 2000);
}

async function connect() {
    if (tryToConnecting) return;
    tryToConnecting = true;
    // disconnect first
    disconnect();
    // script websocket
    socket = new WebSocket("ws://" + location.host + "/ws");
    socket.onmessage = (event) => {
        onMessage(JSON.parse(event.data));
    };
    socket.onopen = () => {
        tryToConnecting = false;
    };
}

function disconnect() {
    if (socket) socket.close();
    Object.values(aircrafts).forEach((aircraft) => {
        aircraft.lines.forEach((line)=>{
            map.removeLayer(line);
        });
        map.removeLayer(aircraft.marker);
    })
    aircrafts = {};
    gID("info").innerHTML = "";
    showError('');
}

function randomColor() {
    return "#000000".replace(/0/g, function () {
        return (~~(Math.random() * 16)).toString(16);
    });
}

function onMessage(msg) {
    const aircraft = msg.aircraft;
    // info client
    gID("info").innerHTML = "seed: " + msg.info.seed + " | total aircrafts: " + msg.info.total_aircrafts;
    // no aircraft in msg -> msg config ?
    if (!aircraft || !aircraft.id || aircraft.id.length === 0) {
        return;
    }
    // aircraft id & pos
    const aircraftID = aircraft.id;
    const pos = aircraft.pos;
    // aircraft info
    let aircraftInfo = aircrafts[aircraftID];
    // aircraft not in dictionary ?
    if (!aircraftInfo) {
        // aircraft is finish ?
        if (aircraft.is_finish) {
            // ignore message
            return;
        }
        // create marker aircraft
        const marker = L.marker(L.latLng(pos.latitude, pos.longitude));
        map.addLayer(marker);
        // create aircraft info
        aircraftInfo = {
            aircraft: aircraft,
            marker: marker,
            color: randomColor(),
            lines: [],
        }
    } else {
        // aircraft is finish ?
        if (aircraft.is_finish) {
            if (aircraftInfo.marker) map.removeLayer(aircraftInfo.marker);
            delete aircrafts[aircraftID];
            return;
        }
        // aircraft position has change ?
        if (pos.longitude !== aircraftInfo.aircraft.pos.longitude || pos.latitude !== aircraftInfo.aircraft.pos.latitude) {
            // move aircraft
            aircraftInfo.marker.setLatLng(L.latLng(pos.latitude, pos.longitude));
            const diffLong =  Math.abs( aircraftInfo.aircraft.pos.longitude-pos.longitude);
            const diffLat =  Math.abs( aircraftInfo.aircraft.pos.latitude-pos.latitude);
            if(diffLat > 5 || diffLong > 5) {
                console.log('diff long', diffLong)
                console.log('diff lat', diffLat)
            }
            // line
            const line = L.polyline([[aircraftInfo.aircraft.pos.latitude, aircraftInfo.aircraft.pos.longitude], [pos.latitude, pos.longitude]], {color: aircraftInfo.color});
            line.addTo(map);
            aircraftInfo.lines.push(line);
            if(aircraftInfo.lines.length > lengthMaxLines) {
                const firstLine = aircraftInfo.lines.shift();
                map.removeLayer(firstLine);
            }

        }
        aircraftInfo.aircraft = aircraft;
    }
    // set aircraft info
    aircrafts[aircraftID] = aircraftInfo;
}