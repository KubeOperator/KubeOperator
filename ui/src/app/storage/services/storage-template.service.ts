import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import {StorageTemplate} from '../models/storage-template';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class StorageTemplateService {
  private storageTemplateUrl = 'api/v1/template/';

  constructor(private http: HttpClient) {
  }

  listStorageTemplates(): Observable<StorageTemplate[]> {
    return this.http.get<StorageTemplate[]>(this.storageTemplateUrl);
  }
}
