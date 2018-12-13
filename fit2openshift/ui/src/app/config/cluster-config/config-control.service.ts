import {Injectable} from '@angular/core';
import {ConfigBase} from './config-base';
import {FormControl, FormGroup} from '@angular/forms';

@Injectable()
export class ConfigControlService {

  constructor() {
  }

  toFormGroup(configs: ConfigBase<any>[]) {
    let group: any = {};

    configs.forEach(config => {
      group[config.name] = new FormControl(config.value);
    });
    return new FormGroup(group);
  }
}
