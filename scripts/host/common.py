#!/usr/bin/env python3

import os
import subprocess
import xml.etree.ElementTree as ET


def project_path():
    return os.popen("git rev-parse --show-toplevel --show-superproject-working-tree").read().strip()


def run_collect_output(cmd):
    res = subprocess.run(cmd, stdout=subprocess.PIPE)
    return res.stdout.decode('utf-8').strip()


def run_in_bg(cmd: str, wd: str):
    return subprocess.Popen(cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                            cwd=os.path.join(project_path(), wd))


def run_shell(cmd):
    res = subprocess.run(cmd, stdout=subprocess.PIPE, shell=True)
    return res.stdout.decode('utf-8').strip()

def addresses_from_manifest(manifest_file: str) -> "list[str]":
    tree = ET.parse(manifest_file)
    root = tree.getroot()
    addresses = []
    for child in root:
        # print(child.tag)
        if child.tag.endswith("node"):
            component_id = child.attrib["component_id"]
            # print(component_id)
            node_name = component_id.split("+")[-1]
            location = component_id.split("+")[1]
            address = f'{node_name}.{location}'
            # print(address)
            addresses.append(address)
    return addresses

PROJECT_PATH = project_path()

# MSSERVERS = ["1106", "1136", "0917", "1117", "1121", "1114"]
# SERVERS = [f'ms{s}.utah.cloudlab.us' for s in MSSERVERS]
# MSSERVERS = ["18", "29", "28", "24", "19", "30"]
# MSSERVERS = [f"0111{s}" for s in MSSERVERS]
# MSSERVERS = ["011121", "011124", "011119"]
# SERVERS = [f"c220g2-{s}.wisc.cloudlab.us" for s in MSSERVERS]
SERVERS = addresses_from_manifest(f'{PROJECT_PATH}/manifest.xml')
