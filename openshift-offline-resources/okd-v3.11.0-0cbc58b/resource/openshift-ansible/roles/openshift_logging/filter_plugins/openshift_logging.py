'''
 Openshift Logging class that provides useful filters used in Logging
'''

import random
import re


def es_storage(os_logging_facts, dc_name, pvc_claim, root='elasticsearch'):
    '''Return a hash with the desired storage for the given ES instance'''
    deploy_config = os_logging_facts[root]['deploymentconfigs'].get(dc_name)
    if deploy_config:
        storage = deploy_config['volumes']['elasticsearch-storage']
        if storage.get('hostPath'):
            return dict(kind='hostpath', path=storage.get('hostPath').get('path'))
    if len(pvc_claim.strip()) > 0:
        return dict(kind='pvc', pvc_claim=pvc_claim)
    return dict(kind='emptydir')


def min_cpu(left, right):
    '''Return the minimum cpu value of the two values given'''
    message = "Unable to evaluate {} cpu value is specified correctly '{}'. Exp whole, decimal or int followed by M"
    pattern = re.compile(r"^(\d*\.?\d*)([Mm])?$")
    millis_per_core = 1000
    if not right:
        return left
    m_left = pattern.match(left)
    if not m_left:
        raise RuntimeError(message.format("left", left))
    m_right = pattern.match(right)
    if not m_right:
        raise RuntimeError(message.format("right", right))
    left_value = float(m_left.group(1))
    right_value = float(m_right.group(1))
    if m_left.group(2) not in ["M", "m"]:
        left_value = left_value * millis_per_core
    if m_right.group(2) not in ["M", "m"]:
        right_value = right_value * millis_per_core
    response = left
    if left_value != min(left_value, right_value):
        response = right
    return response


def walk(source, path, default, delimiter='.'):
    '''Walk the sourch hash given the path and return the value or default if not found'''
    if not isinstance(source, dict):
        raise RuntimeError('The source is not a walkable dict: {} path: {}'.format(source, path))
    keys = path.split(delimiter)
    max_depth = len(keys)
    cur_depth = 0
    while cur_depth < max_depth:
        if keys[cur_depth] in source:
            source = source[keys[cur_depth]]
            cur_depth = cur_depth + 1
        else:
            return default
    return source


def random_word(source_alpha, length):
    ''' Returns a random word given the source of characters to pick from and resulting length '''
    return ''.join(random.choice(source_alpha) for i in range(length))


def entry_from_named_pair(register_pairs, key):
    ''' Returns the entry in key given results provided by register_pairs '''
    results = register_pairs.get("results")
    if results is None:
        raise RuntimeError("The dict argument does not have a 'results' entry. "
                           "Must not have been created using 'register' in a loop")
    for result in results:
        item = result.get("item")
        if item is not None:
            name = item.get("name")
            if name == key:
                return result["content"]
    raise RuntimeError("There was no entry found in the dict that had an item with a name that matched {}".format(key))


def entry_from_name_value_pair(key_value_dict, key, key_label='name', value_label='value'):
    ''' Returns the entry in key given results provided by register_pairs '''
    for key_value in key_value_dict:
        name = key_value.get(key_label)
        if name == key:
            return key_value[value_label]
    # pylint: disable=line-too-long, too-few-format-args
    raise RuntimeError("There was no entry found in the dict that had an item with a name that matched {}:{}".format(key_label).format(key))


def serviceaccount_name(qualified_sa):
    ''' Returns the simple name from a fully qualified name '''
    return qualified_sa.split(":")[-1]


def serviceaccount_namespace(qualified_sa, default=None):
    ''' Returns the namespace from a fully qualified name '''
    seg = qualified_sa.split(":")
    if len(seg) > 1:
        return seg[-2]
    if default:
        return default
    return seg[-1]


def flatten_dict(data, parent_key=None):
    """ This filter plugin will flatten a dict and its sublists into a single dict
    """
    if not isinstance(data, dict):
        raise RuntimeError("flatten_dict failed, expects to flatten a dict")

    merged = dict()

    for key in data:
        if parent_key is not None:
            insert_key = '.'.join((parent_key, key))
        else:
            insert_key = key

        if isinstance(data[key], dict):
            merged.update(flatten_dict(data[key], insert_key))
        else:
            merged[insert_key] = data[key]

    return merged


# pylint: disable=too-few-public-methods
class FilterModule(object):
    ''' OpenShift Logging Filters '''

    # pylint: disable=no-self-use, too-few-public-methods
    def filters(self):
        ''' Returns the names of the filters provided by this class '''
        return {
            'random_word': random_word,
            'entry_from_named_pair': entry_from_named_pair,
            'entry_from_name_value_pair': entry_from_name_value_pair,
            'min_cpu': min_cpu,
            'es_storage': es_storage,
            'serviceaccount_name': serviceaccount_name,
            'serviceaccount_namespace': serviceaccount_namespace,
            'walk': walk,
            "flatten_dict": flatten_dict
        }
