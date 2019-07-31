import {Config} from '../package/package';

export class Region {
  id: string;
  name: string;
  vars: {} = {};
  date_created: string;
  template: string;
  comment: string;
  cloud_region: string;
}

export class CloudTemplate {
  id: string;
  name: string;
  meta: CloudTemplateMeta;
}

export class CloudTemplateMeta {
  name: string;
  meta: Meta;
}

export class Meta {
  configs: Config[] = [];
}
