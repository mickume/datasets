#!/bin/bash

# Example: ./bin/fandom.sh harrypotter fandom_harrypotter

base_dir='./datasets'

# check if values are provided
if [ $# -eq 0 ]; then
    echo "Please provide one or more values as arguments."
    exit 1
fi

data_dir="$base_dir/$1"
fandom="$1"
repo="$2"

# run python script to scrape fandom
python fandomwiki/fandom_wiki.py "$data_dir" "$fandom"

# run wikiextractor with necessary options
wikiextractor "$data_dir/$fandom.xml" --no-templates -l --json -o "$data_dir/$fandom"

# run python script to convert json to text
python fandomwiki/fandom_wiki_markdown.py "$data_dir/$fandom/" "$data_dir/raw/"

# cleanup raw files 
dsc "$data_dir"

# create a dataset and upload it to Huggingface
python create_dataset.py --path "$data_dir/data/" --repo "$repo"

# cleanup temp files
rm -rf "$data_dir/tmp"
rm -rf "$data_dir/$fandom"
rm -f "$data_dir/$fandom.xml"
