#!/bin/bash

# Example: ao3.sh "Harry Potter" ao3_harrypotter ao3_harrypotter 

base_dir='./datasets'
org='mickume'
pages=3

# check if values are provided
if [ $# -eq 0 ]; then
    echo "Please provide one or more values as arguments."
    exit 1
fi

tag="$1"
data_dir="$base_dir/$2"
repo="$3"

# search for story ids
aos "$data_dir/input.txt" "$tag" $pages

# retrieve the stories
aoc $data_dir input.txt

# cleanup raw files 
dsc "$data_dir"

# create a dataset and upload it to Huggingface
python fandomwiki/create_dataset.py --path "$data_dir/data/" --repo "$repo" --user "$org"
