import {Injectable, ViewChild} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Theme} from "./theme";

@Injectable({
    providedIn: 'root'
})
export class ThemeService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/theme';

    get(): Observable<Theme> {
        return this.http.get<Theme>(this.baseUrl);
    }

    update(theme: Theme): Observable<Theme> {
        return this.http.post<Theme>(this.baseUrl, theme);
    }

    setTheme() {
        this.get().subscribe(data => {
            sessionStorage.setItem('theme', JSON.stringify(data));
        });
    }
}
