from common import *
from envs import *


def rsync(src, dst, excludes=[]):
    to_exclude = ['*.pyc', '*.pyo', '*.pyd', '__pycache__', '.idea']
    to_exclude = to_exclude + excludes
    to_exclude = ' '.join([f'--exclude {e}' for e in to_exclude])
    shell_cmd = f'rsync -e "ssh -o StrictHostKeyChecking=no -i {KEYFILE}" -r {to_exclude} {src} {dst}'
    run_shell(shell_cmd)


def main():
    master = SERVERS[0]
    print("copying files to", master)
    path = f'{CLOUDLAB_USER}@{master}:~/'
    rsync(PROJECT_PATH, path, excludes=['architecture_diagram.PNG', 'README.md', 'proxy/target'])
    for s in SERVERS[1:]:
        print("copying files to", s)
        path = f'{CLOUDLAB_USER}@{s}:~/'
        rsync(os.path.join(PROJECT_PATH, "scripts/setup/worker.sh"), path)


if __name__ == '__main__':
    main()
