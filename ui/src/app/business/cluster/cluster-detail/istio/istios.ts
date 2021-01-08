import {BaseModel} from "../../../../shared/class/BaseModel";

export class Istio extends BaseModel {
    name: string;
    version: string;
    describe: string;
    status: string;
    message: string;
    logo: string;
    url: string;
    frame: boolean;
    vars: string;
}

export class IstioHelper {
    cluster_istio: Istio;
    enable: boolean;
    operation: string;
    vars: {} = {};
}