import {BaseModel} from "../../../../shared/class/BaseModel";

export class ClusterTool extends BaseModel {
    name: string;
    version: string;
    describe: string;
    status: string;
    message: string;
    logo: string;
    url: string;
    frame: boolean;
    vars: {} = {};
    conditions: string;
}
