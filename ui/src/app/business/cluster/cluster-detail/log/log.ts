import {BaseModel} from "../../../../shared/class/BaseModel";

export class ClusterLog extends BaseModel {
    message: string;
    type: string;
    status: string;
    startTime: string;
    endTime: string;
}
