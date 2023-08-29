import os
import logging
import argparse
import datasets


LOGGER = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)


def create_dataset(path: str, repo: str) -> None:
    token = os.environ.get('HUGGING_FACE_HUB_TOKEN')
    use_token = True
    if token == None:
        use_token = False

    LOGGER.info(f'Start preparing dataset from {path}')

    ds = datasets.load_dataset(path=path, token=use_token)
    LOGGER.info(f'The dataset is composed of {ds.num_rows} elements.')

    ds.push_to_hub(repo)
    LOGGER.info(f'Uploading dataset finished.')


if __name__ == '__main__':

    parser = argparse.ArgumentParser()
    parser.add_argument('--repo', required=True, help='Name of the Huggingface repo') 
    parser.add_argument('--path', required=True, help='Directory with the data files')
    args = parser.parse_args()

    create_dataset(args.path, args.repo)