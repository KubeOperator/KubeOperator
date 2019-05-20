import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {StorageTemplate} from '../models/storage-template';
import {Observable} from 'rxjs';
import {Storage} from '../models/storage';

@Injectable({
  providedIn: 'root'
})
export class StorageService {

  private storageUrl = 'api/v1/storage/';

  constructor(private http: HttpClient) {
  }
  getStorage(): Observable<Storage[]> {
    return this.http.get<Storage[]>(this.storageUrl);
  }
}

