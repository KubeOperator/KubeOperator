import {V1Node} from "@kubernetes/client-node";
import {CreateNodeRequest} from "../../cluster";
import {BaseModel} from "../../../../shared/class/BaseModel";

export class Node extends BaseModel {
    name: string;
    status: string;
    message: string;
    info: V1Node;
}


export class NodeCreateRequest {
    hosts: CreateNodeRequest[] = [];
    increase: number;
}
