export class Package {
  id: string;
  name: string;
  meta: PackageMeta;
  date_created: string;
}

export class PackageMeta {
  name: string;
  version: string;
  configs: Config[];
  templates: Template[];
}

export class Config {
  name: string;
  alias: string;
  type: string;
  options: Option[];
  value: string;
  default: string;
  require: boolean;
  help_text: string;
}

export class Option {
  name: string;
  alias: string;
}

export class Role {
  name: string;
  meta: RoleMeta;
}

export class RoleMeta {
  hidden: boolean;
  nodes_require: any[];

}

export class Template {
  name: string;
  roles: Role[];
}

