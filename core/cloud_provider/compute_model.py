import logging
import os

import yaml

from fit2ansible.settings import CLOUDS_RESOURCE_DIR

logger = logging.getLogger(__name__)
compute_models = []


def load_compute_model():
    with open((os.path.join(CLOUDS_RESOURCE_DIR, 'compute_model_meta.yml'))) as f:
        logger.info('Load compute model meta')
        compute_models.extend(yaml.load(f))


def get_compute_model_meta(model_name):
    for model in compute_models:
        if model['name'] == model_name:
            return model['meta']
