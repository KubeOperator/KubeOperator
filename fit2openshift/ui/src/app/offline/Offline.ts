import {ConfigBase} from '../config/cluster-config/config-base';

export class Offline {
  id: string;
  name: string;
  location: string;
  branch: string;
  comment: string;
  is_active: boolean;
  config: ConfigBase<string>[];
}
