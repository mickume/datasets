# datasets
My collection of tools, scripts, crawlers etc to create datasets for LLM fine tuning or creation


## Setup for development

### Python virtual environment

```shell
pip install virtualenv

python -m venv venv

source venv/bin/activat
```

#### Install packages

```shell
pip install -r requirements.txt
```

## Examples

### Retrieve the full contents of a Fandom Wiki

Example:

```shell
source venv/bin/activate

./bin/fandom.sh criticalrole "mickume/fandom_criticalrole"
```

### Create a dataset from AO3 stories

Example:

```shell
# search and download stories for tag 'Harry Potter'
./bin/ao3.sh "Harry Potter" small_harrypotter harry_potter_small

# create a tokinized dataset for the model training
./bin/prepare_dataset.sh 'bigscience/bloom-3b' 'mickume/harry_potter_small'

```

### Create a mixed dataset from different fandoms/tags

Example:

First collect a list of stories from different fandoms or tags and combine them into one input file, later used by the crawler:

```shell

# 10 pages/200 stories of 'Hermione Granger + Draco Malfoy'
aos datasets/alt_potterverse/input.txt "Hermione Granger*s*Draco Malfoy" 10

# 10 pages/200 stories of 'Dark Hermione Granger'
aos datasets/alt_potterverse/input.txt "Dark Hermione Granger" 10

# 5 pages/100 stories of 'Harry Potter'
aos datasets/alt_potterverse/input.txt "Harry Potter" 5

```

Now crawl AO3 and create the dataset:

```shell

# retrieve the stories
aoc datasets/alt_potterverse/ input.txt

# cleanup raw files 
dsc datasets/alt_potterverse/

# create a dataset and upload it to Huggingface
python create_dataset.py --path "datasets/alt_potterverse/data/" --repo "mickume/alt_potterverse"
```

All the above in one command:

```shell
clear && \
aos datasets/alt_potterverse/input.txt "Hermione Granger*s*Draco Malfoy" 10 && \
aos datasets/alt_potterverse/input.txt "Dark Hermione Granger" 10 && \
aos datasets/alt_potterverse/input.txt "Harry Potter" 5 && \
aoc datasets/alt_potterverse/ input.txt && \
dsc datasets/alt_potterverse/ && \
python create_dataset.py --path "datasets/alt_potterverse/data/" --repo "mickume/alt_potterverse"
```

Create a tokinzed version of the dataset

```shell
./bin/prepare_dataset.sh 'bigscience/bloom-3b' 'mickume/alt_potterverse'
``
