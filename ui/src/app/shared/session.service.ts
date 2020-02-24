import {Injectable} from '@angular/core';
import {SignInCredential} from './signInCredential';
import {Observable} from 'rxjs';
import {SessionUser} from './session-user';
import {HttpClient} from '@angular/common/http';
import {stringify} from '@angular/compiler/src/util';

const signUrl = '/login';
const authUserUrl = '/api/v1/api-token-auth/';
const getUserUrl = '/api/v1/profile/';
const refreshUrl = '/api/v1/api-token-refresh/';
const userUrl = '/api/v1/users/';
const changePassUrl = '/api/v1/user/{userId}/password';

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

  refreshToken(token: string) {
    return this.http.post<SessionUser>(refreshUrl, {token: token});
  }


  cacheToken(user: SessionUser) {
    localStorage.setItem('current_user', JSON.stringify(user));
  }

  setCacheUser(user: SessionUser) {
    const session = JSON.parse(localStorage.getItem('current_user'));
    session.user = user;
    this.cacheToken(session);
  }

  getCacheUser(): SessionUser {
    let currentUser = null;
    if (localStorage.getItem('current_user') !== null) {
      currentUser = JSON.parse(localStorage.getItem('current_user')).user;
    }
    return currentUser;
  }

  getUser(): Observable<SessionUser> {
    return this.http.get<SessionUser>(getUserUrl);
  }


  changePassword(userId: number, password: string, newPassword: string): Observable<any> {
    const params = {
      password: password,
      new_password: newPassword
    };
    return this.http.post<any>(changePassUrl.replace('{userId}', userId + ''), params);
  }

  clear(): void {
    localStorage.removeItem('current_user');
  }

}
