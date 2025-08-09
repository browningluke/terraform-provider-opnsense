#!/usr/bin/env python3

import base64
import socket
import json
import os
import select
import sys
from time import sleep

QEMU_GA_SOCKET = os.environ.get("QEMU_GA_SOCKET", "/tmp/qemu-virtserialport.sock")


def send_qemu_command(command, socket_path, timeout=5):
    try:
        with socket.socket(socket.AF_UNIX, socket.SOCK_STREAM) as client:
            client.settimeout(timeout)
            client.connect(socket_path)
            client.sendall((json.dumps(command) + "\n").encode("utf-8"))

            response = b""
            while True:
                ready, _, _ = select.select([client], [], [], timeout)
                if not ready:
                    raise TimeoutError(f"No response received within {timeout} seconds")

                chunk = client.recv(4096)
                if not chunk:
                    break
                response += chunk
                if b"\n" in chunk:
                    break
            return response.decode("utf-8")
    except (socket.error, TimeoutError, ConnectionRefusedError) as e:
        return f"Error: {e}"


command_create = {
    "execute": "guest-exec",
    "arguments": {
        "path": "/usr/local/bin/opn-apikey",
        "arg": ["-u", "root", "create"],
        "capture-output": True,
    },
}

output = send_qemu_command(command_create, QEMU_GA_SOCKET).strip()

print(f"Command output: '{output}'")

# {"return": {"pid": 1234}}
pid = json.loads(output).get("return").get("pid")

print(f"Command executed with PID: {pid}")

sleep(5)

command_status = {"execute": "guest-exec-status", "arguments": {"pid": pid}}
output = send_qemu_command(command_status, QEMU_GA_SOCKET).strip()

outdata = json.loads(output).get("return").get("out-data")  # base64 encoded
decoded_outdata = base64.b64decode(outdata).decode("utf-8").strip()

print(decoded_outdata, file=sys.stderr)
