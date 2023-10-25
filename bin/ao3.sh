#!/bin/bash

# Example: ./bin/ao3.sh "Harry Potter" test test

base_dir='./data'
pages=5

# check if values are provided
if [ $# -eq 0 ]; then
    echo "Please provide one or more values as arguments."
    exit 1
fi

tag="$1"
namespace="$2"
repo="$3"

# search for story ids
aos "$base_dir/$namespace/input.txt" "$tag" $pages

# retrieve the stories
aoc $base_dir $namespace input.txt

# cleanup raw files 
dsc $base_dir $namespace input.txt

# create a dataset and upload it to Huggingface
python create_dataset.py --path "$base_dir/$namespace/data/" --repo "$repo"
