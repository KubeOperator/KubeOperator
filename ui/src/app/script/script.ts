export class Script {
  id: string;
  name: string;
  type: string;
  content: string;
  vars: string;
  date_created: string;

  constructor() {
    this.type = 'shell';
  }
}
