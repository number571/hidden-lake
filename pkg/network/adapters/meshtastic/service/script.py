import json
import asyncio
import base64
import threading
import os
import signal
import warnings
import meshtastic.serial_interface

from ipaddress import ip_address
from fastapi import FastAPI, HTTPException, BackgroundTasks, Body
from meshtastic import portnums_pb2
from pubsub import pub

try:
    DEV_PATH = {{devPath}} # example: "/dev/ttyUSB1"
    if DEV_PATH == "":
        raise
except Exception:
    DEV_PATH = None # default

try:
    SRV_ADDR = {{srvAddr}} # example: "127.0.0.1:8080"
    if SRV_ADDR == "":
        raise
except Exception:
    SRV_ADDR = "127.0.0.1:0" # default

# Suppress all deprecation warnings globally
warnings.filterwarnings("ignore", category=DeprecationWarning)

app = FastAPI(title="Meshtastic Serial Binary HTTP Gateway")

mutex = threading.Lock()

received_binary_messages = []
interface = None
is_connected = False

class MeshtasticConnectionError(Exception):
    """Exception raised when a Meshtastic hardware device fails to respond."""
    pass

def on_receive_packet(packet, interface):
    """Callback for processing all incoming packets from the Mesh network."""
    global received_binary_messages
    try:
        if 'decoded' in packet:
            decoded = packet['decoded']
            portnum = decoded.get('portnum')
            channel = packet.get('channel', 0)
            
            if portnum == 'PRIVATE_APP' or 'payload' in decoded:
                payload_bytes = decoded.get('payload', b'')
                
                if payload_bytes:
                    base64_str = base64.b64encode(payload_bytes).decode('utf-8')
                    message_data = {"channel": channel, "message": base64_str}

                    with mutex:
                        received_binary_messages.append(message_data)
    except Exception:
        pass

@app.on_event("startup")
async def startup_event():
    global interface
    global is_connected

    pub.subscribe(on_receive_packet, "meshtastic.receive")

    try:
        interface = meshtastic.serial_interface.SerialInterface(DEV_PATH)
        if interface.myInfo:
            is_connected = True
        else:
            raise MeshtasticConnectionError("Failed connect to node via Serial")
    except Exception:
        pass

@app.on_event("shutdown")
async def shutdown_event():
    global interface
    global is_connected

    if interface:
        interface.close()
        is_connected = False

    os.kill(os.getpid(), signal.SIGTERM)

@app.get("/", summary="Get a list of incoming binary messages")
async def recv_binary_messages():
    """Returns the history of received packets with Base64 strings."""
    global received_binary_messages
    with mutex:
        return_messages = received_binary_messages.copy()
        received_binary_messages = []
    return {"messages": return_messages}

@app.post("/", summary="Send binary data to the mesh network")
async def send_binary_message(json_data: str = Body(...)):
    """Takes Base64, decodes it into bytes and sends it to Mesh."""
    payload = json.loads(json_data)
    
    global interface
    global is_connected

    if not is_connected:
        if interface:
            try:
                interface.close()
            except Exception:
                pass # Ignore errors during forced cleanup
        interface = None
        try:
            interface = meshtastic.serial_interface.SerialInterface(devPath=DEV_PATH)
            is_connected = True
        except Exception as e:
            raise HTTPException(status_code=503, detail="Serial interface is unavailable")
    
    try:
        raw_bytes = base64.b64decode(payload["message"])
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid Base64 format in the message field")
    
    if len(raw_bytes) > 200:
        raise HTTPException(status_code=400, detail="The data size exceeds the Mesh packet limit (~200 bytes)")

    try:
        await asyncio.to_thread(
            interface.sendData,
            data=raw_bytes,
            channelIndex=payload["channel"],
            portNum=portnums_pb2.PortNum.PRIVATE_APP
        )
        return {"status": "success", "bytes_sent": len(raw_bytes)}
    except Exception as e:
        is_connected = False
        raise HTTPException(status_code=500, detail=f"Error sending: {str(e)}")

@app.delete("/", summary="Shutdown HTTP listening service")
async def shutdown(background_tasks: BackgroundTasks):
    """Sends a termination signal to the current process."""
    background_tasks.add_task(shutdown_event)
    return {"status": "success", "message": "Server shutting down..."}

if __name__ == "__main__":
    import uvicorn
    host, _, port = SRV_ADDR.rpartition(":")
    uvicorn.run(app, host=host.strip("[]"), port=int(port))
