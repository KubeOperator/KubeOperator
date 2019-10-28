import {Storage} from '../cluster/storage';

export class NfsStorage extends Storage {
  vars = {
    'allow_ip': '*'
  };
}
