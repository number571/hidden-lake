import meshtastic.util

def get_all_meshtastic_usb_ports() -> list[str]:
    ports = meshtastic.util.findPorts()
    return ports

if __name__ == "__main__":
    print("Find USB devices API Meshtastic...")
    devices = get_all_meshtastic_usb_ports()
    
    if devices:
        print(f"Found Meshtastic devices: {len(devices)}")
        for idx, port in enumerate(devices, 1):
            print(f" {idx}. Port: {port}")
    else:
        print("Meshtastic devices not found. Check USB-connection.")
