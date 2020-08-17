import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Profile} from './session-user';
import {Observable} from 'rxjs';

const queryKey = 'profile';
const profileUrl = '/api/v1/auth/profile/';

@Injectable({
    providedIn: 'root'
})
export class SessionService {

    constructor(private http: HttpClient) {
    }

    cacheProfile(profile: Profile) {
        localStorage.setItem(queryKey, JSON.stringify(profile));
    }

    getCacheProfile(): Profile {
        const profile = localStorage.getItem(queryKey);
        if (profile !== null) {
            return JSON.parse(profile);
        }
        return null;
    }

    getProfile(): Observable<Profile> {
        return this.http.get<Profile>(profileUrl);
    }

    clear() {
        localStorage.removeItem(queryKey);
    }
}
