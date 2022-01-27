import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Captcha, Profile} from './session-user';
import {Observable} from 'rxjs';
import {LoginCredential} from "../../login/login-credential";

const queryKey = 'profile';

@Injectable({
    providedIn: 'root'
})
export class SessionService {

    sessionUrl = '/api/v1/auth/session';
    codeUrl = '/api/v1/captcha';
    profileUrl = '/api/v1/auth/profile';


    constructor(private http: HttpClient) {
    }

    login(item: LoginCredential): Observable<Profile> {
        return this.http.post<Profile>(this.sessionUrl, item);
    }

    logout(): Observable<any> {
        return this.http.delete<any>(this.sessionUrl);
    }

    getCode(): Observable<Captcha> {
        return this.http.get<Captcha>(this.codeUrl);
    }

    clear() {
        sessionStorage.removeItem(queryKey);
    }

    getProfile(): Observable<Profile> {
        return this.http.get<Profile>(this.sessionUrl);
    }
}
