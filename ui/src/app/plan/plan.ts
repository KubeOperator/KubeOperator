export class Plan {
  id: string;
  name: string;
  vars: {} = {};
  date_created: string;
  region: string;
  zones: string[] = [];
}
