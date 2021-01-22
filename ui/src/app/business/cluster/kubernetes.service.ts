import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {V1NamespaceList} from '@kubernetes/client-node/dist/gen/model/v1NamespaceList';
import {Observable} from 'rxjs';
import {
    NetworkingV1beta1IngressList,
    V1beta1CronJobList,
    V1beta1CSINodeList,
    V1ConfigMapList,
    V1CSINodeList,
    V1DaemonSetList,
    V1DeploymentList, V1EventList,
    V1JobList, V1Namespace,
    V1NodeList,
    V1PersistentVolume,
    V1PersistentVolumeClaimList,
    V1PersistentVolumeList,
    V1PodList,
    V1SecretList,
    V1ServiceList,
    V1StatefulSetList,
    V1StorageClass,
    V1StorageClassList
} from '@kubernetes/client-node';

@Injectable({
    providedIn: 'root'
})

export class KubernetesService {

    proxyUrl = '/proxy/kubernetes/{cluster_name}/{resource_url}';
    limit = 10;
    continueTokenKey = 'continue';

    constructor(private client: HttpClient) {
    }

    namespaceUrl = 'api/v1/namespaces';
    namespaceOpUrl = 'api/v1/namespaces/{namespace}';
    serviceUrl = 'api/v1/services';
    namespaceServiceUrl = 'api/v1/namespaces/{namespace}/services';
    persistentVolumesUrl = 'api/v1/persistentvolumes';
    persistentVolumesOpUrl = 'api/v1/persistentvolumes/{name}';
    storageClassUrl = 'apis/storage.k8s.io/v1/storageclasses';
    storageClassOpUrl = 'apis/storage.k8s.io/v1/storageclasses/{name}';
    persistentVolumeClaimsUrl = 'api/v1/persistentvolumeclaims';
    namespacePersistentVolumeClaimsUrl = 'api/v1/namespaces/{namespace}/deployments';
    deploymentUrl = 'apis/apps/v1/deployments';
    namespaceDeploymentUrl = 'apis/apps/v1/namespaces/{namespace}/deployments';
    daemonSetUrl = 'apis/apps/v1/daemonsets/';
    statefulSetUrl = 'apis/apps/v1/statefulsets/';
    namespaceStatefulSet = 'apis/apps/v1/namespaces/{namespace}/statefulsets/';
    namespaceDaemonSetUrl = 'apis/apps/v1/namespaces/{namespace}/daemonsets/';
    cornJobUrl = 'apis/batch/v1beta1/cronjobs';
    namespaceCornJobUrl = 'apis/batch/v1beta1/namespaces/{namespace}/cronjobs';
    jobUrl = 'apis/batch/v1/jobs';
    namespaceJobUrl = 'apis/batch/v1/namespaces/{namespace}/jobs';
    ingressUrl = 'apis/networking.k8s.io/v1beta1/ingresses';
    namespaceIngressUrl = 'apis/networking.k8s.io/v1beta1/namespaces/{namespace}/ingresses';
    configMapUrl = 'api/v1/configmaps';
    namespaceConfigMapUrl = 'api/v1/namespaces/{namespace}/configmaps';
    secretUrl = 'api/v1/secrets';
    namespaceSecretUrl = 'api/v1/secrets';
    podUrl = 'api/v1/pods';
    namespacePodUrl = 'api/v1/namespaces/{namespace}/pods/';
    nodesUrl = 'api/v1/nodes';
    nodeStatsSummaryUrl = 'apis/metrics.k8s.io/v1beta1/nodes';
    eventByNamespaceUrl = 'api/v1/namespaces/{namespace}/events';
    eventsUrl = 'api/v1/events';

    listNodesUsage(clusterName: string, continueToken?: string): Observable<any> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.nodeStatsSummaryUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<any>(url);
    }


    listNodes(clusterName: string, continueToken?: string): Observable<V1NodeList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.nodesUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1NodeList>(url);
    }

    listNamespaces(clusterName: string): Observable<V1NamespaceList> {
        const url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.namespaceUrl);
        return this.client.get<V1NamespaceList>(url);
    }

    createPersistentVolume(clusterName: string, item: V1PersistentVolume): Observable<V1PersistentVolume> {
        const url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.persistentVolumesUrl);
        return this.client.post<V1PersistentVolume>(url, item);
    }

    listPersistentVolumes(clusterName: string, continueToken?: string): Observable<V1PersistentVolumeList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.persistentVolumesUrl);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1PersistentVolumeList>(url);
    }

    listStorageClass(clusterName: string, continueToken?: string, all?: boolean): Observable<V1StorageClassList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.storageClassUrl);
        if (!all) {
            url += '?limit=' + this.limit;
        }
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        return this.client.get<V1StorageClassList>(url);
    }


    createStorageClass(clusterName: string, item: V1StorageClass): Observable<V1StorageClass> {
        const url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.storageClassUrl);
        return this.client.post<V1StorageClass>(url, item);
    }


    listPersistentVolumeClaims(clusterName: string, continueToken?: string, namespace?: string): Observable<V1PersistentVolumeClaimList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespacePersistentVolumeClaimsUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.persistentVolumeClaimsUrl);
        }
        return this.client.get<V1PersistentVolumeClaimList>(url);
    }

    listDeployment(clusterName: string, continueToken?: string, namespace?: string): Observable<V1DeploymentList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceDeploymentUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.deploymentUrl);
        }

        return this.client.get<V1DeploymentList>(url);
    }

    listDaemonSet(clusterName: string, continueToken?: string, namespace?: string): Observable<V1DaemonSetList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceDaemonSetUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.daemonSetUrl);
        }

        return this.client.get<V1DaemonSetList>(url);
    }

    listStatefulSet(clusterName: string, continueToken?: string, namespace?: string): Observable<V1StatefulSetList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceStatefulSet).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.statefulSetUrl);
        }

        return this.client.get<V1StatefulSetList>(url);
    }

    listJob(clusterName: string, continueToken?: string, namespace?: string): Observable<V1JobList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceJobUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.jobUrl);
        }
        return this.client.get<V1JobList>(url);
    }


    listCornJob(clusterName: string, continueToken?: string, namespace?: string): Observable<V1beta1CronJobList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceCornJobUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.cornJobUrl);
        }
        return this.client.get<V1beta1CronJobList>(url);
    }

    listService(clusterName: string, continueToken?: string, namespace?: string): Observable<V1ServiceList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceServiceUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.serviceUrl);
        }
        return this.client.get<V1ServiceList>(url);
    }

    listIngress(clusterName: string, continueToken?: string, namespace?: string): Observable<NetworkingV1beta1IngressList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceIngressUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.ingressUrl);
        }
        return this.client.get<NetworkingV1beta1IngressList>(url);
    }

    listConfigMap(clusterName: string, continueToken?: string, namespace?: string): Observable<V1ConfigMapList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceConfigMapUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.configMapUrl);
        }
        return this.client.get<V1ConfigMapList>(url);
    }

    listSecret(clusterName: string, continueToken?: string, namespace?: string): Observable<V1SecretList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.namespaceSecretUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.secretUrl);
        }
        return this.client.get<V1SecretList>(url);
    }

    listPod(clusterName: string, continueToken?: string, namespace?: string): Observable<V1PodList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        if (namespace) {
            url = url.replace('{resource_url}', this.namespacePodUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.podUrl);
        }
        return this.client.get<V1PodList>(url);
    }


    listEvents(clusterName: string, continueToken?: string, namespace?: string): Observable<V1EventList> {
        let url = this.proxyUrl.replace('{cluster_name}', clusterName);
        url += '?limit=' + this.limit;
        if (continueToken) {
            url += '&continue=' + continueToken;
        }
        if (namespace) {
            url = url.replace('{resource_url}', this.eventByNamespaceUrl).replace('{namespace}', namespace);
        } else {
            url = url.replace('{resource_url}', this.eventsUrl);
        }
        return this.client.get<V1EventList>(url);
    }

    createNamespace(clusterName: string, item: V1Namespace): Observable<V1Namespace> {
        const url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.namespaceUrl);
        return this.client.post<V1Namespace>(url, item);
    }

    deleteNamespace(clusterName: string, namespace: string): Observable<V1Namespace> {
        const url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.namespaceOpUrl).replace('{namespace}', namespace);
        return this.client.delete<V1Namespace>(url);
    }

    deleteStorageClass(clusterName: string, name: string): Observable<V1StorageClass> {
        const url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.storageClassOpUrl).replace('{name}', name);
        return this.client.delete<V1StorageClass>(url);
    }

    deletePersistentVolume(clusterName: string, name: string): Observable<V1PersistentVolume> {
        const url = this.proxyUrl.replace('{cluster_name}', clusterName).replace('{resource_url}', this.persistentVolumesOpUrl).replace('{name}', name);
        return this.client.delete<V1PersistentVolume>(url);
    }
}
