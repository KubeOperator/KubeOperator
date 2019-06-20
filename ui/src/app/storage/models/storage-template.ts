export class StorageTemplate {
  name: string;
  meta: StorageMeta;
}

export class StorageMeta {
  name: string;
  roles: StorageRole[] = [];
  vars: {};
  options: StorageTemplateOption[] = [];
}

export class StorageRole {
  name: string;
  meta: StorageRoleMeta;
}

export class StorageRoleMeta {
  requires: any[] = [];
}

export class StorageTemplateOption {
  name: string;
  comment: string;
  default: string;
  type: string;
  value: any;
  placeholder: string;
}
