import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {V1NamespaceList} from '@kubernetes/client-node/dist/gen/model/v1NamespaceList';
import {Observable} from 'rxjs';
// @ts-ignore
import {IncomingMessage} from 'http';
import {V1NodeList, V1PersistentVolume, V1PersistentVolumeList} from "@kubernetes/client-node";

@Injectable({
    providedIn: 'root'
})

export class KubernetesService {

    proxyUrl = '/api/v1/proxy/{cluster_name}/{resource_url}';
    limit = 10;
    continueTokenKey = 'continue';

    constructor(private client: HttpClient) {
    }

    namespaceUrl = '/api/v1/namespaces';
    persistentVolumesUrl = '/api/v1/persistentvolumes';
    nodesUrl = '/api/v1/nodes';

    listNodes(clusterName: string, continueToken?: string): Observable<V1NodeList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.nodesUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1NodeList>(url);
    }

    listNamespaces(clusterName: string, continueToken?: string): Observable<V1NamespaceList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.namespaceUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1NamespaceList>(url);
    }


    listPersistentVolumes(clusterName: string, continueToken?: string): Observable<V1PersistentVolumeList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.persistentVolumesUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1PersistentVolumeList>(url);
    }
}
