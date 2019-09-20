import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { BackupStorage } from './backup-storage';

@Injectable({
  providedIn: 'root'
})
export class BackupStorageService {

  private baseURL = '/api/v1/backupStorage/';

  constructor(private http: HttpClient) { }

  listBackupStorage(): Observable<BackupStorage[]> {
    return this.http.get<BackupStorage[]>(this.baseURL);
  }

  deleteBackupStorage(name: string): Observable<BackupStorage> {
    return this.http.delete<BackupStorage>(this.baseURL + name + '/');
  }

  createBackupStorage(item: BackupStorage): Observable<BackupStorage> {
    return this.http.post<BackupStorage>(this.baseURL, item);
  }
}
