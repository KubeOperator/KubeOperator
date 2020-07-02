import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {CreateStorageProvisionerRequest, StorageProvisioner} from './storage-provisioner';

@Injectable({
    providedIn: 'root'
})
export class StorageProvisionerService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/clusters/provisioner/{cluster_name}/';

    list(clusterName: string): Observable<StorageProvisioner[]> {
        return this.http.get<StorageProvisioner[]>(this.baseUrl.replace('{cluster_name}', clusterName));
    }

    create(clusterName: string, item: CreateStorageProvisionerRequest): Observable<StorageProvisioner> {
        return this.http.post<StorageProvisioner>(this.baseUrl.replace('{cluster_name}', clusterName), item);
    }
}


