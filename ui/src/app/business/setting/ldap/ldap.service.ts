import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {System} from '../system/system';

@Injectable({
  providedIn: 'root'
})
export class LdapService {

  constructor(private http: HttpClient) {
  }

  baseUrl = '/api/v1/ldap';

  ldapCreate(item): Observable<System> {
    const itemUrl = `${this.baseUrl}`;
    return this.http.post<System>(itemUrl, item);
  }

  ldapSync(item): Observable<System> {
    const itemUrl = `${this.baseUrl}/sync`;
    return this.http.post<System>(itemUrl, item);
  }
}
