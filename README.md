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

./bin/fandom.sh criticalrole
```