import {BaseModel} from "../../../shared/class/BaseModel";

export class CloudProvider extends BaseModel {
    id: string;
    name: string;
    vars: string;
}