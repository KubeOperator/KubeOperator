import {ConfigBase} from './config-base';

export class SelectConfig extends ConfigBase<string> {
  controlType = 'select';
  options: { key: string, value: string }[] = [];

  constructor(options: {} = {}) {
    super(options);
    this.options = options['options'] || [];
  }

}
