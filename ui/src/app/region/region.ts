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
  meta: CloudMeta;
}

export class CloudMeta {
  name: string;
  region: Meta;
  zone: ZoneMeta;
  plan: PlanMeta;
}

export class PlanMeta {
  models: ModelMeta[] = [];
  configs: Config[] = [];
  infra: Infra;
}


export class ModelMeta {
  name: string;
  meta: InfraModel;
}

export class InfraModel {
  cpu: number;
  memory: number;
  disk: number;
}

export class ZoneMeta {
  network: Meta;
  storage: Meta;
  image: Meta;
  others: Meta;
}

export class Meta {
  configs: Config[] = [];
  vars: {} = {};
}

export class Infra {
  name: string;
  alias: string;
}
