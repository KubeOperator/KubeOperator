export class AuthTemplate {
  id: string;
  name: string;
  meta: AuthMeta;
}


export class AuthMeta {
  name: string;
  config: any;
  options: Option[];
  vars: AuthVar[];
}

export class Option {
  name: string;
  default: any;
  children: Option[];
  type: string;
  value: any;
}

export class AuthVar {
  name: string;
  default: any;
  type: string;
  value: any;
}
