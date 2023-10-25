#!/bin/bash

# Example: ./bin/prepare_dataset.sh "model" "dataset"
# e.g. ./bin/prepare_dataset.sh 'bigscience/bloom-3b' 'mickume/harry_potter_small'

base_dir='./datasets'

# check if values are provided
if [ $# -eq 0 ]; then
    echo "Please provide one or more values as arguments."
    exit 1
fi

model="$1"
dataset="$2"
repo="$2_tk"

echo "Creating a dataset for: $dataset"

# create a tokenized dataset and upload it to Huggingface
python prepare_dataset.py --model "$model" --repo "$repo" --dataset "$dataset"
