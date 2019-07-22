import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Credential} from './credential-list/credential';

@Injectable({
  providedIn: 'root'
})
export class CredentialService {
  private baseURl = '/api/v1/credential/';

  constructor(private http: HttpClient) {
  }

  listCredential(): Observable<Credential[]> {
    return this.http.get<Credential[]>(this.baseURl);
  }

  createCredential(item: Credential): Observable<Credential> {
    return this.http.post<Credential>(this.baseURl, item);
  }

  deleteCredential(name: string): Observable<Credential> {
    return this.http.delete<Credential>(this.baseURl + name + '/');
  }

}
