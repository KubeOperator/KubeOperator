import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {CreateStorageProvisionerRequest, ProvisionerSync, StorageProvisioner} from './storage-provisioner';

@Injectable({
    providedIn: 'root'
})
export class StorageProvisionerService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/clusters/provisioner/{cluster_name}';

    list(clusterName: string): Observable<StorageProvisioner[]> {
        return this.http.get<StorageProvisioner[]>(this.baseUrl.replace('{cluster_name}', clusterName));
    }

    create(clusterName: string, item: CreateStorageProvisionerRequest): Observable<StorageProvisioner> {
        return this.http.post<StorageProvisioner>(this.baseUrl.replace('{cluster_name}', clusterName), item);
    }

    syncList(clusterName: string, hosts: ProvisionerSync[]): Observable<any> {
        const syncUrl = '/api/v1/clusters/provisioner/sync/{cluster_name}';
        const url = syncUrl.replace('{cluster_name}', clusterName) ;
        return this.http.post<ProvisionerSync[]>(url, hosts);
    }

    delete(clusterName: string, item: StorageProvisioner): Observable<any> {
        const deleteUrl = '/api/v1/clusters/provisioner/delete/{cluster_name}';
        const url = deleteUrl.replace('{cluster_name}', clusterName) ;
        return this.http.post<any>(url, {item});
    }

    batch(clusterName: string, items: StorageProvisioner[]): Observable<any> {
        const url = this.baseUrl.replace('{cluster_name}', 'batch/' + clusterName);
        return this.http.post<any>(url, {items, operation: 'delete'});
    }
}


