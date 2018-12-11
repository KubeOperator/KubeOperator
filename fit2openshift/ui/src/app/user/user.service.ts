import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {User} from './user';
import {catchError} from 'rxjs/operators';

const userUrl = '/api/v1/users/';

@Injectable()
export class UserService {

  constructor(private http: HttpClient) {
  }

  listUsers(): Observable<User[]> {
    return this.http.get<User[]>(userUrl).pipe(
      catchError(error => throwError(error))
    );
  }

  createUser(user: User): Observable<User> {
    return this.http.post<User>(userUrl, user).pipe(
      catchError(error => throwError(error))
    );
  }

  deleteUser(userId): Observable<any> {
    return this.http.delete(userUrl).pipe(
      catchError(error => throwError(error))
    );
  }

}
