import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {LoginCredential} from './login-credential';
import {Observable} from 'rxjs';
import {Profile} from '../shared/session-user';

@Injectable({
    providedIn: 'root'
})
export class LoginService {

    loginUrl = '/api/login';

    constructor(private http: HttpClient) {
    }

    login(item: LoginCredential): Observable<Profile> {
        return this.http.post<Profile>(this.loginUrl, item);
    }
}
