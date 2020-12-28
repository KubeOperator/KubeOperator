import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';
import {Region} from '../region/region';
import {IpPool} from '../ip-pool/ip-pool';

export class Zone extends BaseModel {
    id: string;
    name: string;
    vars: string;
    credentialId: string;
    cloudVars: {} = {};
    regionName: string;
    provider: string;
    status: string;
    ipPool: IpPool = new IpPool();
}

export class ZoneCreateRequest extends BaseRequest {
    vars: string;
    regionName: string;
    regionID: string;
    cloudVars: {} = {};
    provider: string;
    credentialId: string;
    ipPoolName: string;
}

export class ZoneUpdateRequest extends BaseRequest {
    vars: string;
    regionID: string;
    cloudVars: {} = {};
    ipPoolName: string;
}

export class CloudZoneRequest extends BaseRequest {
    cloudVars: {} = {};
    datacenter: string;
}

export class CloudZone {
    cluster: string;
    networks: [] = [];
    resourcePools: [] = [];
    datastores: [] = [];
    storages: Storage[] = [];
    securityGroups: [] = [];
    networkList: Network[] = [];
    floatingNetworkList: Network[] = [];
    ipTypes: [] = [];
    imageList: Image[] = [];
    switchs: Switch[] = [];
    templates: string[] = [];
}

export class Switch {
    name: string;
    portgroups: string[];
}


export class CloudTemplate {
    imageName: string;
    guestId: string;
}

export class Storage {
    id: string;
    name: string;
}

export class Network {
    id: string;
    name: string;
    subnetList: Subnet[] = [];
}

export class Subnet {
    id: string;
    name: string;
}

export class Image {
    id: string;
    name: string;
}

export class CloudDatastore {
    name: string;
    capacity: number;
    freeSpace: number;
}





