import {Operation} from '../cluster/cluster';

export class Package {
  id: string;
  name: string;
  meta: PackageMeta;
  date_created: string;
  vars: {};
}

export class Network {
  name: string;
  configs: Config[] = [];
}

export class Components {
  kubernetes: string;
  etcd: string;
  docker: string;
}

export class Image {
  name: string;
  tag: string;
}

export class PackageMeta {
  resource: string;
  version: string;
  components: Components;
  images: Image[];
  apps: App[] = [];
  cluster_infos: ClusterInfo[] = [];
  operations: Operation[] = [];
  vars = {};
}

export class App {
  name: string;
  logo: string;
  url_key: string;
  describe: string;
  display_on: string;
}

export class ClusterInfo {
  name: string;
  key: string;
  value: string = null;
}

export class Config {
  name: string;
  alias: string;
  type: string;
  options: Option[];
  value: string;
  display: boolean;
  default: string;
}

export class Option {
  name: string;
  value: string;
}

export class Role {
  name: string;
  meta: RoleMeta;
}

export class Requires {
  nodes_require: any[];
  volumes_require: Require[];
  device_require: Require[];
}

export class Require {
  name: string;
  verbose: string;
  minimal: number;
  excellent: number;
  unit: string;
  comment: string;
}

export class NodeVars {
  name: string;
  template: string;
  verbose: string;
  comment: string;
  type: string;
  options: any;
  placeholder: string;
  require: boolean;
}

export class RoleMeta {
  hidden: boolean;
  allow_os: Os[];
  requires: Requires;
  node_vars: NodeVars[];


}

export class Os {
  name: string;
  version: string[];
}

export class Template {
  name: string;
  deploy_type: string;
  roles: Role[] = [];
  private_config: Config[] = [];
  portals: Portal[];
  comment: string;
}

export class Portal {
  name: string;
  redirect: string;
}
