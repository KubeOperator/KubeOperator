export class ResultMessage {
  id: string;
  message: string;
  type: string;
  category: string;
}

export class CountSummary {
  successes: number;
  warnings: number;
  errors: number;
}

export class ResultSummary {
  by_category: {};
  totals: CountSummary;
}

export class ContainerResult {
  summary: ResultSummary;
  name: string;
  messages: ResultMessage[] = [];
}

export class PodResult {
  summary: ResultSummary;
  messages: ResultMessage[] = [];
  container_results: ContainerResult[] = [];
}

export class ControllerResult {
  name: string;
  namespace: string;
  type: string;
  pod_result: PodResult;
}

export class NamespaceResult {
  name: string;
  controller_results: ControllerResult[] = [];
  summary: ResultSummary;
}

export class ClusterResult {
  namespace_results: NamespaceResult[] = [];
  summary: ResultSummary;
  grade: string;
  score: number;
}
