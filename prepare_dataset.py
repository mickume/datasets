import os
import logging
import argparse

from pathlib import Path
from typing import Callable, Mapping, Any

from datasets import Dataset, load_dataset
from transformers import AutoTokenizer, PreTrainedTokenizer


LOGGER = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)


def tokenize(element: Mapping, tokenizer: Callable, context_length: int) -> str:
    inputs = tokenizer(element['text'],
                       truncation=True,
                       return_overflowing_tokens=True,
                       return_length=True,
                       max_length=context_length)

    inputs_batch = []
    for length, input_ids in zip(inputs['length'], inputs['input_ids']):
        if length == context_length:  # We drop the last input_ids that are shorter than max_length
            inputs_batch.append(input_ids)

    return {"input_ids": inputs_batch}


def preprocess_dataset(dataset_name: str, repo_name: str, split: str, tokenizer: PreTrainedTokenizer, context_length: int, test_size: float = 0.1, shuffle: bool = True) -> None:
    LOGGER.info(f'Start preparing dataset {dataset_name}')

    dataset = load_dataset(dataset_name, split=split)

    large_text = ''
    for para in dataset:
        large_text += para['text'] + ' '#+ tokenizer.eos_token

    dataset = Dataset.from_dict({'text': [large_text]})

    tokenized_dataset = dataset.map(tokenize, batched=True, fn_kwargs={
                                    'tokenizer': tokenizer, 'context_length': context_length}, remove_columns=dataset.column_names)
    LOGGER.info(
        f'The tokenized dataset is composed of {tokenized_dataset.num_rows} elements, each one composed of {context_length} tokens.')

    tokenized_dataset_dict = tokenized_dataset.train_test_split(
        test_size=test_size, shuffle=shuffle)
    LOGGER.info(
        f'The training dataset is composed of {tokenized_dataset_dict["train"].num_rows} elements, the test dataset is composed of {tokenized_dataset_dict["test"].num_rows} elements.')

    tokenized_dataset_dict.push_to_hub(repo_name)
    LOGGER.info(f'Uploading tokenized dataset finished.')


if __name__ == '__main__':

    parser = argparse.ArgumentParser()
    parser.add_argument('--model', type=str, required=True,
                        help='Name of the base model')
    parser.add_argument('--dataset', type=str, required=True,
                        help='Name of the Huggingface repo with the dataset')
    parser.add_argument('--repo', type=str, required=True,
                        help='Name of the Huggingface repo with the resulting tokenized dataset')
    parser.add_argument('--split', type=str, default='train')
    parser.add_argument('--context_length', type=int, default=2048)
    parser.add_argument('--test_size', type=float, default=0.1)
    args = parser.parse_args()

    tokenizer = AutoTokenizer.from_pretrained(args.model)
    preprocess_dataset(args.dataset, args.repo, args.split,
                       tokenizer, args.context_length, args.test_size, True)
