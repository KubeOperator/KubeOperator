export class Execution {
  id: string;
  date_start: string;
  date_end: string;
  date_created: string;
  log_ws_url: string;
  log_url: string;
  progress_ws_url: string;
  state: string;
  operation: string;
  timedelta: string;
  steps: Step[];
}

export class Step {
  name: string;
  status: string;
  flow: string;
}
