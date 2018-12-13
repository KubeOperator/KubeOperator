import { Injectable } from '@angular/core';
import { HttpClient } from "@angular/common/http";
import {Observable, of} from "rxjs/index";

import { AuthUser } from "./auth-user";


@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private usersUrl = '/api/v1';
  currentUser: AuthUser = null;

  constructor(private http: HttpClient) { }

  getProfile(): Observable<AuthUser> {
    const url =  `${this.usersUrl}/profile`;
    return this.http.get<AuthUser>(url)
  }

  getCurrentUser() {
    return this.currentUser
  }

  log(msg) {
    console.log(msg)
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      this.log(`${operation} failed: ${error.message}`);
      return of(result as T);
    };
  }
}
