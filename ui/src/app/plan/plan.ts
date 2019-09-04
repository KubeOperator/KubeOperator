export class Plan {
  id: string;
  name: string;
  vars: {} = {};
  date_created: string;
  region: string;
  zone: string;
  zones: string[] = [];
  deploy_template: string;
}
