import {BaseModel} from './BaseModel';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Page} from './Page';
import {Batch} from './Batch';

export abstract class BaseModelService<T extends BaseModel> {

    baseUrl = '';

    protected constructor(protected http: HttpClient) {
    }

    list(): Observable<T[]> {
        return this.http.get<T[]>(this.baseUrl);
    }

    page(page, size): Observable<Page<T>> {
        const pageUrl = `${this.baseUrl}/${page}/${size}/`;
        return this.http.get<Page<T>>(pageUrl);
    }

    get(name: string): Observable<T> {
        const itemUrl = `${this.baseUrl}/${name}/`;
        return this.http.get<T>(itemUrl);
    }

    create(item: T): Observable<T> {
        return this.http.post<T>(this.baseUrl, item);
    }

    update(name: string, item: T): Observable<T> {
        const itemUrl = `${this.baseUrl}/${name}/`;
        return this.http.patch<T>(itemUrl, item);
    }

    delete(name: string): Observable<any> {
        const itemUrl = `${this.baseUrl}/${name}/`;
        return this.http.delete<any>(itemUrl);
    }

    batch(method: string, items: T[]): Observable<any> {
        const batchUrl = `${this.baseUrl}/batch/`;
        const b = new Batch<T>(method, items);
        return this.http.post(batchUrl, b);
    }

}

