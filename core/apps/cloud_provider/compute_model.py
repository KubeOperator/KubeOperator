import logging
import os
import yaml
from kubeoperator.settings import CLOUDS_RESOURCE_DIR

__all__ = ["compute_models", "load_compute_model", "get_compute_model_meta"]
logger = logging.getLogger(__name__)
compute_models = []


def load_compute_model():
    with open((os.path.join(CLOUDS_RESOURCE_DIR, 'compute_model_meta.yml'))) as f:
        compute_models.extend(yaml.load(f))


def get_compute_model_meta(model_name):
    for model in compute_models:
        if model['name'] == model_name:
            return model['meta']
