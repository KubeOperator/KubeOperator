import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {StorageNode} from '../models/storage-node';

@Injectable({
  providedIn: 'root'
})
export class StorageNodeService {
  private storageNodeUrl = '/api/v1/storage/{storage_name}/nodes/';

  constructor(private http: HttpClient) {
  }

  listStorageNode(storageName: string): Observable<StorageNode[]> {
    return this.http.get<StorageNode[]>(this.storageNodeUrl.replace('{storage_name}', storageName));
  }

  createStorageNode(storageName: string, item: StorageNode): Observable<StorageNode> {
    return this.http.post<StorageNode>(this.storageNodeUrl.replace('{storage_name}', storageName), item);
  }

}
