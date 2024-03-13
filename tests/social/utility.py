import requests
import logging


def heartbeat(ips, app: str):
    url = compose_url(ips, app, "heartbeat")
    r = requests.get(url)
    logging.debug(f'Heartbeat response code: {r.status_code}')
    logging.debug(f'Heartbeat response text: {r.text}')
    assert (r.status_code == 200)
    assert (r.text.rstrip() == "Heartbeat")


# ips: service_name -> ip
def compose_url(ips, app: str, method: str):
    return compose_url_raw(method, ips[app])


def compose_url_raw(method: str, ip: str = "localhost", port: int = 80):
    url = f'http://{ip}:{port}/{method}'
    return url


def populate_ips_from_args(args):
    # don't abort if a service is not specified
    # used for unit tests
    ips = {service: getattr(args, service, None) for service in
           ["post_storage", "home_timeline", "user_timeline", "social_graph", "compose_post"]}
    return ips

def general_ips_from_args(args):
#     print(args, dir(args), args.__dict__())
    return vars(args)

def create_user_id_from_int(index: int) -> str:
    return f'User{index}'


## Social graph related code

class SocialGraph:
    def __init__(self, number_of_nodes):
        ## Nodes are immutable
        self.nodes = list(range(1, number_of_nodes + 1))
        self.edges = {node: [] for node in self.nodes}
        self.edges_inv = {node: [] for node in self.nodes}

    def add_edge(self, from_node, to_node):
        self.edges[from_node].append(to_node)
        self.edges_inv[to_node].append(from_node)

    def get_nodes(self):
        return self.nodes

    def degree_in(self, node: int) -> int:
        return len(self.edges_inv[node])

    def get_incoming_nodes(self, node: int) -> list:
        return self.edges_inv[node]
    
    def get_outgoing_nodes(self, node: int) -> list:
        return self.edges[node]


def parse_social_graph(social_graph_file: str) -> SocialGraph:
    with open(social_graph_file) as f:
        data = f.readlines()

    data_no_comments = [line for line in data
                        if not line.startswith("%")]

    rows, columns, edges_number = data_no_comments[0].split()
    nodes = int(rows)
    assert (nodes == int(columns))

    edge_lines = data_no_comments[1:]
    assert (int(edges_number) == len(edge_lines))

    social_graph = SocialGraph(nodes)
    for line in edge_lines:
        from_node_s, to_node_s = line.split()
        social_graph.add_edge(int(from_node_s), int(to_node_s))

    return social_graph


## Analyzes the graph and stores a file
def analyze_social_graph(social_graph: SocialGraph, analysis_file: str, post_size: int, number_of_posts_per_user: int):
    nodes = social_graph.get_nodes()

    followers = {}
    total_followers = 0
    for node in nodes:
        node_followers = social_graph.degree_in(node)
        followers[node] = node_followers
        total_followers += node_followers

    ## Not sure if we need the ratios
    follower_ratio = {}
    for node, node_followers in followers.items():
        follower_ratio[node] = float(node_followers) / total_followers

    with open(analysis_file, "w") as f:
        ## Write number of nodes
        f.write(f'Nodes: {len(nodes)}\n')

        ## Total followers
        f.write(f'Total followers: {total_followers}\n')

        f.write(f'Post size (in bytes): {post_size}\n')

        f.write(f'Number of starting posts per user: {number_of_posts_per_user}\n')

        ## Write follower ratio
        sorted_nodes = sorted(list(followers.items()), key=lambda x: x[1], reverse=True)
        f.write(f'Followers (nodes sorted by follower count):\n')
        for node, node_followers in sorted_nodes:
            f.write(f'{node} {node_followers}\n')

    return total_followers
