#!/bin/bash

# Example: ./bin/create_dataset.sh criticalrole fandom_criticalrole

base_dir='./datasets'

# check if values are provided
if [ $# -eq 0 ]; then
    echo "Please provide one or more values as arguments."
    exit 1
fi

data_dir="$base_dir/$1"
dataset="$1"
repo="$2"

echo "Creating a dataset for: $dataset"

# cleanup raw files 
dsc "$data_dir"

# create a dataset and upload it to Huggingface
python create_dataset.py --path "$data_dir/data/" --repo "$repo"
