import argparse
import xml.etree.ElementTree as ET

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--file",
                        help="the manifest file to read from",
                        type=str,
                        default="./manifest.xml")
    args = parser.parse_args()
    return args

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

def main():
    args = parse_arguments()
    addresses = addresses_from_manifest(args.file)
    print(len(addresses), "nodes")
    
    address_string = '", "'.join(addresses)
    address_string = f'"{address_string}"'
    print(address_string)

if __name__ == '__main__':
    main()
