#!/usr/bin/env python3
"""Upgrade the QEMU-booted OPNsense VM to the latest 26.1.x patch level.

The bsd.ac qemu image ships OPNsense 26.1 base. Several DNAT bugs were fixed
in 26.1.1-26.1.3 (see community changelog), causing /api/firewall/d_nat/addRule
to return 500 on the base image. Run opnsense-update to apply the patches
before tests.
"""

import base64
import json
import os
import select
import socket
import sys
import time

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


def guest_exec(path, args, max_wait_s=900):
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
    err = base64.b64decode(data.get("err-data", "")).decode("utf-8", errors="replace")
    code = data.get("exitcode", -1)
    return code, out, err


def ping():
    try:
        resp = send({"execute": "guest-ping"}, timeout=3)
        return resp == '{"return": {}}'
    except (socket.error, OSError):
        return False


def wait_for_agent(max_wait_s=180):
    deadline = time.monotonic() + max_wait_s
    while time.monotonic() < deadline:
        if ping():
            return True
        time.sleep(3)
    return False


def main():
    print("Refreshing pkg repository...", flush=True)
    code, out, err = guest_exec("/usr/sbin/pkg", ["update", "-f"], max_wait_s=180)
    print(f"pkg update exit={code}")
    if out:
        print(out)
    if err:
        print(f"[stderr] {err}", file=sys.stderr)

    print("Upgrading opnsense package (carries the MVC d_nat fixes)...", flush=True)
    # -y auto-confirms, -f reinstalls even if same version, then upgrade
    code, out, err = guest_exec("/usr/sbin/pkg", ["upgrade", "-y", "opnsense"], max_wait_s=600)
    print(f"pkg upgrade opnsense exit={code}")
    if out:
        print(out)
    if err:
        print(f"[stderr] {err}", file=sys.stderr)
    if code != 0:
        # Fall back to opnsense-update with no flags (default is firmware update)
        print("pkg upgrade failed; trying opnsense-update -p (firmware patch)...", flush=True)
        code, out, err = guest_exec("/usr/local/sbin/opnsense-update", ["-p"], max_wait_s=600)
        print(f"opnsense-update -p exit={code}")
        if out:
            print(out)
        if err:
            print(f"[stderr] {err}", file=sys.stderr)
        if code != 0:
            sys.exit(f"both pkg upgrade and opnsense-update failed; last exit={code}")

    print("Restarting configd + lighttpd to pick up the new MVC code...", flush=True)
    guest_exec("/usr/local/sbin/configctl", ["webgui", "restart"], max_wait_s=60)
    time.sleep(5)

    print("Verifying new version...", flush=True)
    code, out, _ = guest_exec("/usr/local/sbin/opnsense-version", [], max_wait_s=20)
    print(f"opnsense-version exit={code} output={out.strip()}")


if __name__ == "__main__":
    main()
