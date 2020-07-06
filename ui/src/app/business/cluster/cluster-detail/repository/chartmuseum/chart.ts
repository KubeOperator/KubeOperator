export class Chart {
    name: string;
    home: string;
    version: string;
    description: string;
    apiVersion: string;
    appVersion: string;
    deprecated: boolean;
    urls: string[] = [];
    keywords: string[] = [];
    sources: string[] = [];
    created: string;
    digest: string;
}

export class ChartMap {
    [key: string]: Chart;
}
