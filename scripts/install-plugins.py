#!/usr/bin/env python3
"""Install the plugins listed in opnsense-plugins.txt on the QEMU-booted VM.

Runs the same script as System > Firmware > Plugins, which also registers the
plugin in config.xml. It refuses to install plugins on an outdated core, so
run upgrade-opnsense.py first.
"""

import base64
import json
import os
import pathlib
import select
import socket
import sys
import time

QEMU_GA_SOCKET = os.environ.get("QEMU_GA_SOCKET", "/tmp/qemu-virtserialport.sock")
PLUGINS_FILE = pathlib.Path(__file__).parent / "opnsense-plugins.txt"


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


def guest_exec(path, args, max_wait_s=600):
    resp = send({
        "execute": "guest-exec",
        "arguments": {"path": path, "arg": args, "capture-output": True},
    })
    pid = json.loads(resp)["return"]["pid"]
    deadline = time.monotonic() + max_wait_s
    data = {}
    while time.monotonic() < deadline:
        time.sleep(5)
        try:
            data = json.loads(send({
                "execute": "guest-exec-status",
                "arguments": {"pid": pid},
            })).get("return", {})
        except json.JSONDecodeError:
            continue
        if data.get("exited"):
            break
    out = base64.b64decode(data.get("out-data", "")).decode("utf-8", errors="replace")
    return data.get("exitcode", -1), out


def main():
    plugins = [
        line.strip()
        for line in PLUGINS_FILE.read_text().splitlines()
        if line.strip() and not line.startswith("#")
    ]
    failed = []
    for plugin in plugins:
        print(f"Installing {plugin}...", flush=True)
        _, out = guest_exec("/usr/local/opnsense/scripts/firmware/install.sh", [plugin])
        print(out)
        # install.sh exits 0 even when it refuses, so check the package itself.
        code, _ = guest_exec("/usr/sbin/pkg", ["info", "-e", plugin], max_wait_s=30)
        if code != 0:
            failed.append(plugin)
    if failed:
        sys.exit(f"plugins not installed: {', '.join(failed)}")


if __name__ == "__main__":
    main()
