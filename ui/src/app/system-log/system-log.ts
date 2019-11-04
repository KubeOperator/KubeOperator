export class SystemLogPager {
  items: SystemLog[] = [];
  total: number;
}

export class SystemLog {
  name: string;
  timestamp: string;
  level: string;
  filename: string;
  funcName: string;
  lineno: string;
  message: string;
  host_ip: string;
  exc_text: string;
}
