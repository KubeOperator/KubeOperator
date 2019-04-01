export class ConfigBase<T> {
  value: T;
  name: string;
  comment: string;
  required: boolean;
  controlType: string;

  constructor(options: {
    value?: T,
    name?: string,
    comment?: string,
    required?: boolean,
    controlType?: string
  } = {}) {
    this.value = options.value;
    this.name = options.name || '';
    this.comment = options.comment || '';
    this.required = !!options.required;
    this.controlType = options.controlType || '';
  }
}
