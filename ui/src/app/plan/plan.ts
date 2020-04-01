export class Plan {
  id: string;
  name: string;
  vars: {} = {};
  date_created: string;
  region: string;
  zone: string;
  zones: string[] = [];
  provider: string;
  deploy_template: string;
  item_id: [];
}

export class ComputeModel {
  name: string;
  meta: ComputeModelMeta;
}

export class ComputeModelMeta {
  cpu: number;
  memory: number;
  disk: number;
}
