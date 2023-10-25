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

### Create a dataset from AO3 stories

Example:

```shell
# search and download stories for tag 'Harry Potter'
./bin/ao3.sh "Harry Potter" small_harrypotter harry_potter_small

```

### Create a mixed dataset from different fandoms/tags

Example:

First collect a list of stories from different fandoms or tags and combine them into one input file, later used by the crawler:

```shell

# 10 pages/200 stories with tag 'Hermione Granger + Draco Malfoy'
aos data/alt_potterverse/input.txt "Hermione Granger*s*Draco Malfoy" 10

# 10 pages/200 stories with tag 'Dark Hermione Granger'
aos data/alt_potterverse/input.txt "Dark Hermione Granger" 10

# 5 pages/100 stories with tag 'Harry Potter'
aos data/alt_potterverse/input.txt "Harry Potter" 5

```

Now retrieve the text and create the dataset:

```shell

# retrieve the stories
aoc data alt_potterverse input.txt

# cleanup raw files 
dsc data alt_potterverse input.txt

# create a dataset and upload it to Huggingface
python create_dataset.py --path "datas/alt_potterverse/data/" --repo "mickume/alt_potterverse"
```

All the above in one command:

```shell
namespace="alt_potterverse"

aos data/$namespace/input.txt "Hermione Granger*s*Draco Malfoy" 10 && \
aos data/$namespace/input.txt "Dark Hermione Granger" 10 && \
aos data/$namespace/input.txt "Harry Potter" 5 && \
aoc data $namespace input.txt && \
dsc data $namespace input.txt && \
python create_dataset.py --path "data/$namespace/data/" --repo "mickume/$namespace"
```

### More examples

```shell
./bin/ao3.sh "World of Warcraft" wow wow
```

```shell
namespace="dnd_drow"

aos data/$namespace/input.txt "Original Drow Character%28s%29 %28Dungeons *a* Dragons%29" 10 && \
aos data/$namespace/input.txt "Drow (Dungeons *a* Dragons)" 10 && \
aos data/$namespace/input.txt "Original Dungeons *a* Dragons Character(s)" 5 && \
aos data/$namespace/input.txt "Dungeons%20*a*%20Dragons%20(Roleplaying%20Game)" 5 && \
aoc data $namespace input.txt && \
dsc data $namespace input.txt && \
python create_dataset.py --path "data/$namespace/data/" --repo "mickume/$namespace"
```