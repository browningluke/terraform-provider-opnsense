#!/usr/bin/env python3
"""Send a graceful power-off command to the QEMU-booted OPNsense VM."""

import json
import os
import select
import socket

QEMU_GA_SOCKET = os.environ.get("QEMU_GA_SOCKET", "/tmp/qemu-virtserialport.sock")


def send(command, timeout=10):
    with socket.socket(socket.AF_UNIX, socket.SOCK_STREAM) as client:
        client.settimeout(timeout)
        client.connect(QEMU_GA_SOCKET)
        client.sendall((json.dumps(command) + "\n").encode())
        chunks = b""
        while True:
            ready, _, _ = select.select([client], [], [], timeout)
            if not ready:
                break
            chunk = client.recv(8192)
            if not chunk:
                break
            chunks += chunk
            if b"\n" in chunk:
                break
        return chunks.decode().strip()


def main():
    print("Sending power-off command via QEMU guest agent...", flush=True)
    try:
        resp = send({
            "execute": "guest-exec",
            "arguments": {
                "path": "/sbin/shutdown",
                "arg": ["-p", "now"],
                "capture-output": False,
            },
        })
        print(f"Response: {resp}", flush=True)
    except Exception as exc:
        print(f"Warning: could not send shutdown command: {exc}", flush=True)


if __name__ == "__main__":
    main()
