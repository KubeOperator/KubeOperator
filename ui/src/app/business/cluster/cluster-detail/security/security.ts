import {BaseModel} from "../../../../shared/class/BaseModel";

export class CisTask extends BaseModel {
    id: string;
    startTime: string;
    endTime: string;
    message: string;
    status: string;
    results: CisTaskResult[] = [];
}

export class CisTaskResult {
    number: string;
    desc: string;
    remediation: string;
    status: string;
}
