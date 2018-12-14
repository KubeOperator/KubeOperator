import {Injectable} from '@angular/core';
import {SignInCredential} from './signInCredential';
import {Observable} from 'rxjs';
import {SessionUser} from './session-user';
import {HttpClient} from '@angular/common/http';

const signUrl = '/login';
const authUserUrl = '/api/v1/api-token-auth/';
const getUserUrl = '/api/v1/profile/';

@Injectable({
  providedIn: 'root'
})
export class SessionService {

  constructor(private http: HttpClient) {
  }

  authUser(signInCredential: SignInCredential): Observable<SessionUser> {
    const credential = {
      username: signInCredential.principal,
      password: signInCredential.password
    };
    return this.http.post<SessionUser>(authUserUrl, credential);
  }

  getCacheUser(): SessionUser {
    return JSON.parse(localStorage.getItem('current_user')).user;
  }

  getUser(): Observable<SessionUser> {
    return this.http.get<SessionUser>(getUserUrl);
  }

  clear(): void {
    localStorage.removeItem('current_user');
  }

}
