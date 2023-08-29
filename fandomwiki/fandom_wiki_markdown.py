import json
import os
import re
import argparse
import errno
import html

from tqdm import tqdm
from markdownify import MarkdownConverter

from fandom_wiki_util import CustomConverter, safe_open_w


parser = argparse.ArgumentParser()
parser.add_argument('input_dir', help='Directory with wiki json files')
parser.add_argument('output', help='Txt file output')
args = parser.parse_args()

directories = os.listdir(args.input_dir)
counter = 0


def md(htxt, **options):
    return bytes(CustomConverter(**options).convert(html.unescape(htxt)), 'utf-8')


for directory in tqdm(directories):
    for filename in tqdm(os.listdir(os.path.join(args.input_dir, directory)), desc="Processing "+directory):
        if not filename.startswith('wiki'):
            continue

        path = os.path.join(os.path.join(args.input_dir, directory), filename)
        with open(path, 'r') as fin:
            for line in fin:
                data = json.loads(line)

                if data['text'] == "":
                    continue
                else:
                    fname = data['id'] + ".txt"
                    opath = os.path.join(args.output, fname)

                    title = data['title']+"\n\n"

                    with safe_open_w(opath) as fout:
                        fout.write(bytes(title, 'utf-8'))
                        fout.write(md(data['text']))
