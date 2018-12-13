export class Play {
  id: string;
  name: string;
  project: string;
  pattern: string;
  vars: object;
  tasks: [object];
  roles: [object];
}

export class Git {
  repo: string;
  branch: string;
}

export class Playbook {
  id: string;
  name: string;
  alias: string;
  type: string;
  git: Git;
  project: string;
  plays: [Play];
  is_periodic: boolean;
  interval: string;
  crontab: string;
  is_active: boolean;
  comment: string;
  created_by: string;
  date_created: string;

  constructor() {
    this.git = new Git();
    this.type = 'git';
  }
}

export class PlaybookExecution {
  id: string;
  num: number;
  state: string;
  timedelta: number;
  log_url: string;
  log_ws_url: string;
  result_summary: object;
  date_created: string;
  date_start: string;
  date_end: string;
  comment: string;
}
